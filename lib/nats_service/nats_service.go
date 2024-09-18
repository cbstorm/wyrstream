package nats_service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/cbstorm/wyrstream/lib/configs"
	"github.com/cbstorm/wyrstream/lib/exceptions"
	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/cbstorm/wyrstream/lib/utils"
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

func (ns *NATS_Service) Start(subscriber_id string) {
	s := ns.subscribers[subscriber_id]
	if s == nil {
		return
	}
	go s.Start(ns.nats_client, ns.queue_group)
}

func (ns *NATS_Service) StartAllSubscriber() error {
	for _, s := range ns.subscribers {
		go s.Start(ns.nats_client, ns.queue_group)
	}
	return nil
}

func (ns *NATS_Service) StopSubscribe(subscriber_id string) error {
	s := ns.subscribers[subscriber_id]
	if s == nil {
		return nil
	}
	return s.Stop()
}

type RequestOpts struct {
	timeout time.Duration
	output  interface{}
}

type RequestOptFunc func(*RequestOpts)

func WithTimeout(t time.Duration) RequestOptFunc {
	return func(ro *RequestOpts) {
		ro.timeout = t
	}
}

func WithOutput(o interface{}) RequestOptFunc {
	return func(ro *RequestOpts) {
		ro.output = o
	}
}

func (ns *NATS_Service) Request(subj NATS_Subject, data interface{}, opts ...RequestOptFunc) (interface{}, error) {
	var bytes_data []byte
	switch v := data.(type) {
	case []byte:
		bytes_data = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		bytes_data = b
	}
	o := &RequestOpts{}
	for _, v := range opts {
		v(o)
	}
	if o.timeout == 0 {
		o.timeout = DEFAULT_TIMEOUT
	}
	msg, err := ns.nats_client.Request(string(subj), bytes_data, o.timeout)
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
	if o.output != nil {
		if err := utils.Cast(res.Data, o.output); err != nil {
			return nil, err
		}
	}
	return res.Data, nil
}
