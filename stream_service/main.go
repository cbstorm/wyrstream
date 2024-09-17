package main

import (
	"github.com/cbstorm/wyrstream/lib/configs"
	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/cbstorm/wyrstream/lib/nats_service"
	local_configs "github.com/cbstorm/wyrstream/stream_service/configs"
	"github.com/cbstorm/wyrstream/stream_service/server"
)

func main() {
	logg := logger.NewLogger("STREAM_SERVICE")
	if err := local_configs.GetConfig().Load(); err != nil {
		logg.Fatal("Could not load local configs with err: %v", err)
	}
	if err := configs.GetConfig().Load(); err != nil {
		logg.Fatal("Could not load configs with err: %v", err)
	}
	if err := nats_service.GetNATSService().Connect(); err != nil {
		logg.Fatal("Could not connect to NATS server with err: %v", err)
	}
	logg.Info("Stream service started.")
	server.GetServer().Init().Listen()
}
