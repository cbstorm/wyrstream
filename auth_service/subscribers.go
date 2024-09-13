package main

import (
	"log"
	"time"

	natsservice "github.com/cbstorm/wyrstream/lib/nats_service"
	"github.com/nats-io/nats.go"
)

var _ = natsservice.GetNATSService().AddSubcriber(natsservice.NewSubscriber("test_1", func(m *nats.Msg) ([]byte, error) {
	log.Println(string(m.Data))
	time.Sleep(5 * time.Second)
	return []byte("response_1"), nil
}))
