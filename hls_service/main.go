package main

import (
	"github.com/cbstorm/wyrstream/lib/configs"
	"github.com/cbstorm/wyrstream/lib/database"
	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/cbstorm/wyrstream/lib/nats_service"
	"github.com/cbstorm/wyrstream/lib/utils"
)

func main() {
	logg := logger.NewLogger("HLS_SERVICE")
	if err := utils.AssertDir(PUBLIC_DIR + "/"); err != nil {
		logg.Fatal("Could not assert dir [./public] with err: %v", err)
	}
	// DB
	db := database.GetDatabase()
	if err := db.LoadConfig(configs.GetConfig()); err != nil {
		logg.Fatal("Could not load database config with err: %v", err)
	}
	if err := db.Connect(); err != nil {
		logg.Fatal("Could not connect to database with err:  %v", err)
	}
	// NATS
	n := nats_service.GetNATSService()
	if err := n.LoadConfig(configs.GetConfig()); err != nil {
		logg.Fatal("Could not load NATS core config with err: %v", err)
	}
	if err := n.Connect(); err != nil {
		logg.Fatal("Could not connect to NATS server with err: %v", err)
	}
	if err := n.StartAllSubscriber(); err != nil {
		logg.Fatal("Could not start subscribers with error: %v", err)
	}
	// HLS HTTP server
	s := GetHttpServer()
	if err := s.LoadConfig(configs.GetConfig()); err != nil {
		logg.Error("Could not load HLS HTTP server config with err: %v", err)
	}
	if err := s.Init().Listen(); err != nil {
		logg.Fatal("Could not start http server with err: %v", err)
	}
}
