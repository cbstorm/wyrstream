package server

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/cbstorm/wyrstream/lib/nats_service"
	"github.com/cbstorm/wyrstream/stream_service/configs"
	srt "github.com/datarhei/gosrt"
)

var instance *Server
var instance_sync sync.Once

func GetServer() *Server {
	if instance == nil {
		instance_sync.Do(func() {
			instance = &Server{}
		})
	}
	return instance
}

type Server struct {
	addr       string
	app        string
	token      string
	passphrase string
	logtopics  string
	server     *srt.Server
	channels   map[string]srt.PubSub
	lock       sync.RWMutex
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown() {
	s.server.Shutdown()
}

func (s *Server) Init() *Server {
	cfg := configs.GetConfig()
	s.addr = cfg.ADDR
	s.app = "/live"
	s.channels = make(map[string]srt.PubSub)
	return s
}

func (s *Server) Listen() {
	if len(s.addr) == 0 {
		fmt.Fprintf(os.Stderr, "Provide a listen address with -addr\n")
		os.Exit(1)
	}

	config := srt.DefaultConfig()

	if len(s.logtopics) != 0 {
		config.Logger = srt.NewLogger(strings.Split(s.logtopics, ","))
	}

	config.KMPreAnnounce = 200
	config.KMRefreshRate = 10000

	s.server = &srt.Server{
		Addr:            s.addr,
		HandleConnect:   s.handleConnect,
		HandlePublish:   s.handlePublish,
		HandleSubscribe: s.handleSubscribe,
		Config:          &config,
	}

	fmt.Fprintf(os.Stderr, "Listening on %s\n", s.addr)

	go func() {
		if config.Logger == nil {
			return
		}

		for m := range config.Logger.Listen() {
			fmt.Fprintf(os.Stderr, "%#08x %s (in %s:%d)\n%s \n", m.SocketId, m.Topic, m.File, m.Line, m.Message)
		}
	}()

	go func() {
		if err := s.ListenAndServe(); err != nil && err != srt.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "SRT Server: %s\n", err)
			os.Exit(2)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	s.Shutdown()

	if config.Logger != nil {
		config.Logger.Close()
	}
}

func (s *Server) log(who, action, path, message string, client net.Addr) {
	fmt.Fprintf(os.Stderr, "%-10s %10s %s (%s) %s\n", who, action, path, client, message)
}

func (s *Server) handleConnect(req srt.ConnRequest) srt.ConnType {
	var mode srt.ConnType = srt.SUBSCRIBE
	client := req.RemoteAddr()
	streamId := req.StreamId()
	path := streamId

	if strings.HasPrefix(streamId, "publish:") {
		mode = srt.PUBLISH
		path = strings.TrimPrefix(streamId, "publish:")
	} else if strings.HasPrefix(streamId, "subscribe:") {
		path = strings.TrimPrefix(streamId, "subscribe:")
	}

	u, err := url.Parse(path)
	if err != nil {
		return srt.REJECT
	}

	if req.IsEncrypted() {
		if err := req.SetPassphrase(s.passphrase); err != nil {
			s.log("CONNECT", "FORBIDDEN", u.Path, err.Error(), client)
			return srt.REJECT
		}
	}

	token := u.Query().Get("token")
	if len(s.token) != 0 && s.token != token {
		s.log("CONNECT", "FORBIDDEN", u.Path, "invalid token ("+token+")", client)
		return srt.REJECT
	}

	if !strings.HasPrefix(u.Path, s.app) {
		s.log("CONNECT", "FORBIDDEN", u.Path, "invalid app", client)
		return srt.REJECT
	}

	if len(strings.TrimPrefix(u.Path, s.app)) == 0 {
		s.log("CONNECT", "INVALID", u.Path, "stream name not provided", client)
		return srt.REJECT
	}

	channel := u.Path
	key := u.Query().Get("key")

	if key == "" {
		s.log("CONNECT", "UNAUTHORIZE", u.Path, "", client)
		return srt.REJECT
	}

	if err := s.onConnect(channel, key, mode); err != nil {
		return srt.REJECT
	}

	s.lock.RLock()
	pubsub := s.channels[channel]
	s.lock.RUnlock()

	if mode == srt.PUBLISH && pubsub != nil {
		s.log("CONNECT", "CONFLICT", channel, "already publishing", client)
		return srt.REJECT
	}

	if mode == srt.SUBSCRIBE && pubsub == nil {
		s.log("CONNECT", "NOTFOUND", channel, "not publishing", client)
		return srt.REJECT
	}

	return mode
}

func (s *Server) handlePublish(conn srt.Conn) {
	client := conn.RemoteAddr()
	if client == nil {
		conn.Close()
		return
	}
	streamId := conn.StreamId()
	path := strings.TrimPrefix(streamId, "publish:")

	u, err := url.Parse(path)
	if err != nil {
		conn.Close()
		return
	}
	channel := u.Path
	key := u.Query().Get("key")
	if key == "" {
		conn.Close()
		return
	}
	if err := s.onPublish(channel, key); err != nil {
		conn.Close()
		return
	}
	log.Printf("[PUBLISH] channel: %s, key: %s", channel, key)

	s.lock.Lock()
	pubsub := s.channels[channel]
	if pubsub != nil {
		s.log("PUBLISH", "CONFLICT", channel, "already publishing", client)
		conn.Close()
		return
	}
	// Init new pubsub
	pubsub = srt.NewPubSub(srt.PubSubConfig{
		Logger: s.server.Config.Logger,
	})
	s.channels[channel] = pubsub
	s.lock.Unlock()

	s.log("PUBLISH", "START", channel, "publishing", client)

	pubsub.Publish(conn)

	s.lock.Lock()
	delete(s.channels, channel)
	s.lock.Unlock()

	s.log("PUBLISH", "STOP", channel, "", client)

	conn.Close()
}

func (s *Server) handleSubscribe(conn srt.Conn) {
	client := conn.RemoteAddr()
	if client == nil {
		conn.Close()
		return
	}

	streamId := conn.StreamId()
	path := strings.TrimPrefix(streamId, "subscribe:")
	u, err := url.Parse(path)
	if err != nil {
		conn.Close()
		return
	}
	channel := u.Path
	key := u.Query().Get("key")
	if key == "" {
		conn.Close()
		return
	}
	if err := s.onSubscribe(channel, key); err != nil {
		conn.Close()
		return
	}
	log.Printf("[SUBSCRIBE] channel: %s, key: %s", channel, key)
	s.log("SUBSCRIBE", "START", channel, "", client)

	s.lock.RLock()
	pubsub := s.channels[channel]
	s.lock.RUnlock()

	if pubsub == nil {
		s.log("SUBSCRIBE", "NOTFOUND", channel, "not publishing", client)
		conn.Close()
		return
	}

	pubsub.Subscribe(conn)

	s.log("SUBSCRIBE", "STOP", channel, "", client)

	conn.Close()
}

func (s *Server) onConnect(stream_id, key string, mode srt.ConnType) error {
	input := dtos.CheckStreamKeyInput{
		StreamId: stream_id,
		Key:      key,
	}
	if mode == srt.PUBLISH {
		res, err := nats_service.GetNATSService().Request(nats_service.AUTH_STREAM_CHECK_PUBLISH_KEY, input)
		if err != nil {
			return err
		}
		if result, ok := res.(*dtos.CheckStreamKeyResponse); !ok || !result.Ok {
			return fmt.Errorf("publish key invalid")
		}
	}
	if mode == srt.SUBSCRIBE {
		res, err := nats_service.GetNATSService().Request(nats_service.AUTH_STREAM_CHECK_SUBSCRIBE_KEY, input)
		if err != nil {
			return err
		}
		if result, ok := res.(*dtos.CheckStreamKeyResponse); !ok || !result.Ok {
			return fmt.Errorf("subscribe key invalid")
		}
	}
	return nil
}

func (s *Server) onPublish(stream_id, publish_key string) error {
	input := dtos.CheckStreamKeyInput{
		StreamId: stream_id,
		Key:      publish_key,
	}
	res, err := nats_service.GetNATSService().Request(nats_service.AUTH_STREAM_CHECK_PUBLISH_KEY, input)
	if err != nil {
		return err
	}
	if result, ok := res.(*dtos.CheckStreamKeyResponse); !ok || !result.Ok {
		return fmt.Errorf("publish key invalid")
	}
	return nil
}

func (s *Server) onPublishStop(stream_id string) error {
	return nil
}

func (s *Server) onSubscribe(stream_id, subscribe_key string) error {
	input := dtos.CheckStreamKeyInput{
		StreamId: stream_id,
		Key:      subscribe_key,
	}
	res, err := nats_service.GetNATSService().Request(nats_service.AUTH_STREAM_CHECK_SUBSCRIBE_KEY, input)
	if err != nil {
		return err
	}
	if result, ok := res.(*dtos.CheckStreamKeyResponse); !ok || !result.Ok {
		return fmt.Errorf("subscribe key invalid")
	}
	return nil
}

func (s *Server) onSubscribeStop(stream_id string) error {
	return nil
}
