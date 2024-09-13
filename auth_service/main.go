package main

import (
	"os"
	"os/signal"

	"github.com/cbstorm/wyrstream/lib/configs"
	"github.com/cbstorm/wyrstream/lib/logger"
	nats_service "github.com/cbstorm/wyrstream/lib/nats_service"
)

func main() {
	logg := logger.NewLogger("AUTH_SERVICE")
	if err := configs.GetConfig().Load(); err != nil {
		logg.Fatal("Could not load configuration with error:%v ", err)
	}
	if err := nats_service.GetNATSService().Connect(); err != nil {
		logg.Fatal("Could not connect to NATS server with error: %v", err)
	}
	if err := nats_service.GetNATSService().StartAllSubscriber(); err != nil {
		logg.Fatal("Could not start subscribers with error: %v", err)
	}
	logg.Info("Auth service started.")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logg.Info("Control service shutting down...")
	os.Exit(0)
}
