package main

import (
	"os"
	"os/signal"

	"github.com/cbstorm/wyrstream/lib/configs"
	"github.com/cbstorm/wyrstream/lib/database"
	"github.com/cbstorm/wyrstream/lib/logger"
	natsservice "github.com/cbstorm/wyrstream/lib/nats_service"
)

func main() {
	logg := logger.NewLogger("CONTROL_SVC")
	if err := configs.GetConfig().Load(); err != nil {
		logg.Fatal("Could not load config with err: %v", err)
	}
	if err := database.GetDatabase().Connect(); err != nil {
		logg.Fatal("Could not connect to database %v", err)
	}
	if err := natsservice.GetNATSService().Connect(); err != nil {
		logg.Fatal("Could not connect to NATS server with err: %v", err)
	}
	logg.Info("Control service started.")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logg.Info("Control service shutting down...")
	os.Exit(0)
}
