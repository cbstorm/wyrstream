package main

import (
	"os"
	"os/signal"

	"github.com/cbstorm/wyrstream/lib/configs"
	"github.com/cbstorm/wyrstream/lib/database"
	"github.com/cbstorm/wyrstream/lib/logger"
	nats_service "github.com/cbstorm/wyrstream/lib/nats_service"
	"github.com/cbstorm/wyrstream/lib/redis_service"
)

func main() {
	logg := logger.NewLogger("AUTH_SERVICE")
	if err := configs.GetConfig().Load(); err != nil {
		logg.Fatal("Could not load configuration with error:%v ", err)
	}
	// DB
	db := database.GetDatabase()
	if err := db.LoadConfig(configs.GetConfig()); err != nil {
		logg.Fatal("Could not load database config with err: %v", err)
	}
	if err := db.Connect(); err != nil {
		logg.Fatal("Could not connect to database with err: %v", err)
	}
	// Redis
	rd := redis_service.GetRedisService()
	if err := rd.LoadConfig(configs.GetConfig()); err != nil {
		logg.Fatal("Could not load redis config with err: %v", err)
	}
	if err := rd.Connect(); err != nil {
		logg.Fatal("Could not connect to redis server with err: %v", err)
	}
	// NATS
	n := nats_service.GetNATSService()
	if err := n.LoadConfig(configs.GetConfig()); err != nil {
		logg.Fatal("Could not load NATS core config with err: %v", err)
	}
	if err := n.Connect(); err != nil {
		logg.Fatal("Could not connect to NATS server with error: %v", err)
	}
	if err := n.StartAllSubscriber(); err != nil {
		logg.Fatal("Could not start subscribers with error: %v", err)
	}
	// Graceful shutdown
	logg.Info("Auth service started.")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logg.Info("Control service shutting down...")
	os.Exit(0)
}
