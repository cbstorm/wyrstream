package main

import (
	"github.com/cbstorm/wyrstream/lib/configs"
	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/cbstorm/wyrstream/lib/nats_service"
	"github.com/cbstorm/wyrstream/lib/utils"
)

func main() {
	logg := logger.NewLogger("HLS_SERVICE")
	if err := utils.AssertDir("public/"); err != nil {
		logg.Fatal("Could not assert dir [./public] with err: %v", err)
	}
	if err := GetConfig().Load(); err != nil {
		logg.Fatal("Could not load local configs with err: %v", err)
	}
	if err := configs.GetConfig().Load(); err != nil {
		logg.Fatal("Could not load configs with err: %v", err)
	}
	if err := nats_service.GetNATSService().Connect(); err != nil {
		logg.Fatal("Could not connect to NATS server with err: %v", err)
	}
	if err := nats_service.GetNATSService().StartAllSubscriber(); err != nil {
		logg.Fatal("Could not start subscribers with error: %v", err)
	}
	if err := GetHttpServer().Init().Listen(); err != nil {
		logg.Fatal("Could not start http server with err: %v", err)
	}
}
