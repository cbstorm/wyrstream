package main

import (
	"log"
	"time"

	nats_service "github.com/cbstorm/wyrstream/lib/nats_service"
	"github.com/nats-io/nats.go"
)

var _ = nats_service.GetNATSService().AddSubcriber(nats_service.NewSubscriber("test_1", func(m *nats.Msg) ([]byte, error) {
	log.Println(string(m.Data))
	time.Sleep(5 * time.Second)
	return []byte("response_1"), nil
}))
