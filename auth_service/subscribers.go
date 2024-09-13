package main

import (
	natsservice "github.com/cbstorm/wyrstream/lib/nats_service"
	"github.com/nats-io/nats.go"
)

var _ = natsservice.GetNATSService().AddSubcriber(natsservice.NewSubscriber("test_1", func(m *nats.Msg) ([]byte, error) {
	return []byte("response_1"), nil
}))
