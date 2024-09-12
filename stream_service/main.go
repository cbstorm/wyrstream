package main

import (
	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/cbstorm/wyrstream/stream_service/configs"
	"github.com/cbstorm/wyrstream/stream_service/server"
)

func main() {
	logg := logger.NewLogger("STREAM_SERVICE")
	if err := configs.GetConfig().Load(); err != nil {
		logg.Fatal("Could not load config with err %v", err)
	}
	logg.Info("Stream service started.")
	server.GetServer().Init().Listen()
}
