package nats_service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/cbstorm/wyrstream/lib/configs"
	"github.com/cbstorm/wyrstream/lib/exceptions"
	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/nats-io/nats.go"
)

var instance *NATS_Service
var instance_sync sync.Once

func GetNATSService() *NATS_Service {
	if instance == nil {
		instance_sync.Do(func() {
			instance = &NATS_Service{
				logger: logger.NewLogger("NATS_SERVICE"),
				config: configs.GetConfig(),
			}
		})
	}
	return instance
}

type NATS_Service struct {
	nats_client *nats.Conn
	logger      *logger.Logger
	queue_group string
	subscribers map[string]*Subscriber
	mu          sync.RWMutex
	config      *configs.Config
}

func (ns *NATS_Service) Connect() error {
	ns.queue_group = ns.config.NATS_CORE_QUEUE_GROUP
	nc_connection_string := fmt.Sprintf("nats://%s:%s@%s:%d", ns.config.NATS_CORE_USERNAME, ns.config.NATS_CORE_PASSWORD, ns.config.NATS_CORE_HOST, ns.config.NATS_CORE_PORT)
	nc, err := nats.Connect(nc_connection_string,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(5),
		nats.ReconnectWait(time.Second*5),
	)
	if err != nil {
		return err
	}
	ns.nats_client = nc
	if err = ns.verifyConnection(); err != nil {
		return err
	}
	ns.logger.Info("Connect to NATS server at %s:%d successfully!", ns.config.NATS_CORE_HOST, ns.config.NATS_CORE_PORT)
	return nil
}

func (ns *NATS_Service) verifyConnection() error {
	version := ns.nats_client.ConnectedServerVersion()
	if version == "" {
		ns.logger.Fatal("%s", "Connected to NATs server failed")
	}
	ns.logger.Info("NATS_Service max_payload: %d", ns.nats_client.MaxPayload())
	ns.logger.Info("Connected to NATs server with version: %s", version)
	return nil
}

func (ns *NATS_Service) GetNC() *nats.Conn {
	return ns.nats_client
}

func (ns *NATS_Service) AddSubcriber(s *Subscriber) bool {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	if ns.subscribers == nil {
		ns.subscribers = make(map[string]*Subscriber)
	}
	ns.subscribers[s.id] = s
	return true
}

func (ns *NATS_Service) StartAllSubscriber() error {
	for _, s := range ns.subscribers {
		go s.Start(ns.nats_client, ns.queue_group)
	}
	return nil
}

type RequestOpts struct {
	timeout time.Duration
}

type RequestOptFunc func(*RequestOpts)

func WithTimeout(t time.Duration) RequestOptFunc {
	return func(ro *RequestOpts) {
		ro.timeout = t
	}
}

func (ns *NATS_Service) Request(subj string, data []byte, opts ...RequestOptFunc) (interface{}, error) {
	o := &RequestOpts{}
	for _, v := range opts {
		v(o)
	}
	if o.timeout == 0 {
		o.timeout = DEFAULT_TIMEOUT
	}
	msg, err := ns.nats_client.Request(subj, data, o.timeout)
	if err != nil {
		return nil, err
	}
	res := &ResponseMessage{}
	if err := json.Unmarshal(msg.Data, res); err != nil {
		return nil, err
	}
	if res.Error != nil {
		err, _ := res.Error.(map[string]interface{})
		e := exceptions.NewException(err["name"].(string)).SetMessage(err["message"].(string)).SetStatus(int(err["status"].(float64)))
		return nil, e
	}
	return res.Data, nil
}
