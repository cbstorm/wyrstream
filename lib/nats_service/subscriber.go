package nats_service

import (
	"encoding/json"
	"fmt"

	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type SubcriberOpts struct {
	concurrency uint32
}

type SubcriberOptFunc func(*SubcriberOpts)

func WithConcurrency(concurrency uint32) SubcriberOptFunc {
	return func(so *SubcriberOpts) {
		so.concurrency = concurrency
	}
}

type SubscriberStatus uint8

const (
	START SubscriberStatus = 0
	STOP  SubscriberStatus = 1
)
const (
	DEFAULT_CC = 10 // Default subscribe concurrency
)

type Subscriber struct {
	id, subject          string
	concurrency          uint32
	handler              func(IRequestMessage) (interface{}, error)
	status               SubscriberStatus
	logger               *logger.Logger
	subscription         *nats.Subscription
	subscription_channel chan *nats.Msg // Close when unsubscribe
}

func NewSubscriber(subj NATS_Subject, handler func(IRequestMessage) (interface{}, error), opts ...SubcriberOptFunc) *Subscriber {
	op := &SubcriberOpts{}
	for _, v := range opts {
		v(op)
	}
	if op.concurrency == 0 {
		op.concurrency = 10
	}
	id, _ := uuid.NewV6()
	id_str := id.String()
	return &Subscriber{
		id:          id_str,
		subject:     string(subj),
		handler:     handler,
		concurrency: op.concurrency,
		status:      STOP,
		logger:      logger.NewLogger(fmt.Sprintf("SUBSCRIBER - <%s> - <%s>", subj, id_str)),
	}
}

func (s *Subscriber) GetId() string {
	return s.id
}

func (s *Subscriber) SetSubject(subj string) *Subscriber {
	s.subject = subj
	return s
}

func (s *Subscriber) SetHandler(handler func(IRequestMessage) (interface{}, error)) *Subscriber {
	s.handler = handler
	return s
}

func (s *Subscriber) Start(nc *nats.Conn, queue_group string) error {
	sub_ch := make(chan *nats.Msg, 64)
	limit_ch := make(chan bool, s.concurrency)
	sub, err := nc.ChanQueueSubscribe(s.subject, queue_group, sub_ch)

	s.subscription = sub
	s.subscription_channel = sub_ch

	defer sub.Drain()

	if err != nil {
		s.logger.Error("Could not start subscibe with error: %v", err)
		return err
	}
	s.status = START
	s.logger.Info("Started")
	for v := range sub_ch {
		go func(msg *nats.Msg) {
			limit_ch <- true
			defer func() {
				<-limit_ch
			}()
			m := &RequestMessage{data: msg.Data}
			result, err := s.handler(m)
			res := &ResponseMessage{Data: result, Error: err}
			r, err := json.Marshal(res)
			if err != nil {
				s.logger.Error("%v", err)
				return
			}
			if err := msg.Respond(r); err != nil {
				s.logger.Error("%v", err)
				return
			}
		}(v)
	}
	return nil
}

func (s *Subscriber) Stop() error {
	if s.subscription != nil {
		if err := s.subscription.Unsubscribe(); err != nil {
			s.logger.Error("Could not unsubscribe with err: %v", err)
			return err
		}
		close(s.subscription_channel)
		s.logger.Info("Unsubscribed")
	}
	return nil
}
