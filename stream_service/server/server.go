package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/cbstorm/wyrstream/lib/nats_service"
	"github.com/cbstorm/wyrstream/lib/redis_service"
	srt "github.com/datarhei/gosrt"
)

type IStreamServerConfig interface {
	LoadStreamServerConfig() error
	SERVER_ADDRESS() string
	SERVER_PUBLIC_URL() string
}

var instance *Server
var instance_sync sync.Once

func GetServer() *Server {
	if instance == nil {
		instance_sync.Do(func() {
			instance = &Server{
				connections: make(map[string][]srt.Conn),
			}
		})
	}
	return instance
}

type Server struct {
	addr          string
	public_url    string
	app           string
	token         string
	passphrase    string
	logtopics     string
	server        *srt.Server
	channels      map[string]srt.PubSub
	lock          sync.RWMutex
	config        IStreamServerConfig
	connections   map[string][]srt.Conn
	publish_count uint

	redis_service *redis_service.RedisService
	nats_service  *nats_service.NATS_Service
}

func (s *Server) LoadConfig(config IStreamServerConfig) error {
	if err := config.LoadStreamServerConfig(); err != nil {
		return err
	}
	s.config = config
	return nil
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown() {
	s.server.Shutdown()
}

func (s *Server) Init() *Server {
	s.addr = s.config.SERVER_ADDRESS()
	s.public_url = s.config.SERVER_PUBLIC_URL()
	s.app = "/live/"
	s.channels = make(map[string]srt.PubSub)

	// deps
	s.redis_service = redis_service.GetRedisService()
	s.nats_service = nats_service.GetNATSService()
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
	log.Printf("REMOTE_ADDR: [%s], STREAMID: [%s] CONNECT", client, streamId)

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
	stream_id := strings.TrimPrefix(channel, s.app)
	key := u.Query().Get("key")
	if key == "" {
		s.log("CONNECT", "UNAUTHORIZE", u.Path, "", client)
		return srt.REJECT
	}

	if err := s.onConnect(stream_id, key, mode); err != nil {
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
	// Remove app prefix
	stream_id := strings.TrimPrefix(channel, s.app)
	key := u.Query().Get("key")
	if key == "" {
		conn.Close()
		return
	}
	if err := s.onPublish(stream_id, key); err != nil {
		s.log("PUBLISH", "ON_PUBLISH_ERROR", channel, err.Error(), client)
		conn.Close()
		return
	}

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

	// Emit START event
	if err := s.onPublishStart(stream_id); err != nil {
		s.log("PUBLISH", "EMIT_START_EVENT", channel, err.Error(), client)
		conn.Close()
		return
	}
	// Append connection by stream_id
	s.connections[stream_id] = append(s.connections[stream_id], conn)
	s.log("PUBLISH", "START", channel, "publishing", client)

	pubsub.Publish(conn)

	s.lock.Lock()
	delete(s.channels, channel)
	s.lock.Unlock()

	s.log("PUBLISH", "STOP", channel, "", client)

	conn.Close()
	// Emit stop event
	s.onPublishStop(stream_id)
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
	stream_id := strings.TrimPrefix(channel, s.app)
	key := u.Query().Get("key")
	if key == "" {
		conn.Close()
		return
	}
	if err := s.onSubscribe(stream_id, key); err != nil {
		s.log("SUBSCRIBE", "ON_SUBSCRIBE_ERROR", channel, "", client)
		conn.Close()
		return
	}
	// Append connection by stream_id
	s.connections[stream_id] = append(s.connections[stream_id], conn)
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
	if mode == srt.PUBLISH {
		return s.onPublish(stream_id, key)
	}
	if mode == srt.SUBSCRIBE {
		return s.onSubscribe(stream_id, key)
	}
	return nil
}

func (s *Server) onPublish(stream_id, publish_key string) error {
	input := &dtos.CheckStreamKeyInput{
		StreamServer: s.public_url,
		StreamId:     stream_id,
		Key:          publish_key,
	}
	result := &dtos.CheckStreamKeyResponse{}
	_, err := s.nats_service.Request(nats_service.AUTH_STREAM_CHECK_PUBLISH_KEY, input, nats_service.WithOutput(result))
	if err != nil {
		return err
	}
	if !result.Ok {
		return errors.New("publish key invalid")
	}
	return nil
}

func (s *Server) onPublishStart(stream_id string) error {
	s.lock.Lock()
	s.publish_count++
	s.lock.Unlock()
	input := &dtos.HLSPublishStartInput{
		StreamId:        stream_id,
		StreamServer:    s.public_url,
		StreamServerApp: s.app,
	}
	_, err := s.nats_service.Request(nats_service.HLS_PUBLISH_START, input)
	if err != nil {
		return err
	}
	return err
}

func (s *Server) onPublishStop(stream_id string) {
	s.lock.Lock()
	s.publish_count--
	s.lock.Unlock()
	input := &dtos.HLSPublishStopInput{
		StreamId: stream_id,
	}
	go func() {
		_, err := s.nats_service.Request(nats_service.HLS_PUBLISH_STOP.Concat(input.StreamId), input)
		if err != nil {
			log.Printf("Could not emit publish stop event of stream %s with err: %v", stream_id, err)
		}
	}()
}

func (s *Server) onSubscribe(stream_id, subscribe_key string) error {
	input := dtos.CheckStreamKeyInput{
		StreamServer: s.public_url,
		StreamId:     stream_id,
		Key:          subscribe_key,
	}
	result := &dtos.CheckStreamKeyResponse{}
	_, err := s.nats_service.Request(nats_service.AUTH_STREAM_CHECK_SUBSCRIBE_KEY, input, nats_service.WithOutput(result))
	if err != nil {
		return err
	}
	if !result.Ok {
		return fmt.Errorf("subscribe key invalid")
	}
	return nil
}

func (s *Server) TerminateStream(stream_id string) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	conns := s.connections[stream_id]
	for _, c := range conns {
		if err := c.Close(); err != nil {
			log.Printf("Could not close conn with error: %v", err)
		}
	}
}

func (s *Server) ConfirmHealth() {
	d := 5 * time.Second
	ticker := time.NewTicker(d)
	t := func() error {
		if err := s.redis_service.SetWithTtl(redis_service.REDIS_KEY_STREAM_SERVER_HEALTH.Concat(s.public_url), s.publish_count, d); err != nil {
			log.Printf("Could not confirm health with err: %v", err)
			return err
		}
		return nil
	}
	if err := t(); err != nil {
		return
	}
	go func() {
		defer ticker.Stop()
		for {
			<-ticker.C
			if err := t(); err != nil {
				break
			}
		}
	}()
}

func (s *Server) BuildStreamSubscribeUrl(stream_id string, key string) string {
	return fmt.Sprintf("srt://%s?streamid=%s%s?key=%s", s.public_url, s.app, stream_id, key)
}
