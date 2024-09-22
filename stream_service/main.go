package main

import (
	"github.com/cbstorm/wyrstream/lib/configs"
	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/cbstorm/wyrstream/lib/nats_service"
	"github.com/cbstorm/wyrstream/lib/redis_service"
	"github.com/cbstorm/wyrstream/stream_service/server"
)

func main() {
	logg := logger.NewLogger("STREAM_SERVICE")
	// NATS
	n := nats_service.GetNATSService()
	if err := n.LoadConfig(configs.GetConfig()); err != nil {
		logg.Fatal("Could not load NATS core config with err: %v", err)
	}
	if err := n.Connect(); err != nil {
		logg.Fatal("Could not connect to NATS server with err: %v", err)
	}
	// Redis
	rd := redis_service.GetRedisService()
	if err := rd.LoadConfig(configs.GetConfig()); err != nil {
		logg.Fatal("Could not load redis config with err: %v", err)
	}
	if err := rd.Connect(); err != nil {
		logg.Fatal("Could not connect to redis server with err: %v", err)
	}
	logg.Info("Stream service started.")
	// Stream server
	s := server.GetServer()
	if err := s.LoadConfig(configs.GetConfig()); err != nil {
		logg.Fatal("Could not load stream server config with err: %v", err)
	}
	s.Init().Listen()
}
