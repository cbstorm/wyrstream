package main

import (
	"github.com/cbstorm/wyrstream/lib/configs"
	"github.com/cbstorm/wyrstream/lib/database"
	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/cbstorm/wyrstream/lib/minio_service"
	"github.com/cbstorm/wyrstream/lib/nats_service"
	"github.com/cbstorm/wyrstream/lib/utils"
)

func main() {
	logg := logger.NewLogger("HLS_SERVICE")
	if err := utils.AssertDir(PUBLIC_DIR + "/"); err != nil {
		logg.Fatal("Could not assert directory [public] due to an error: %v", err)
	}
	// DB
	db := database.GetDatabase()
	if err := db.LoadConfig(configs.GetConfig()); err != nil {
		logg.Fatal("Could not load database configuration due to an error: %v", err)
	}
	if err := db.Connect(); err != nil {
		logg.Fatal("Could not connect to the database due to an error:  %v", err)
	}
	// NATS
	n := nats_service.GetNATSService()
	if err := n.LoadConfig(configs.GetConfig()); err != nil {
		logg.Fatal("Could not load NATS core config due to an error: %v", err)
	}
	if err := n.Connect(); err != nil {
		logg.Fatal("Could not connect to the NATS server due to an error: %v", err)
	}
	if err := n.StartAllSubscriber(); err != nil {
		logg.Fatal("Could not start subscribers due to an error: %v", err)
	}
	// MinIO
	m := minio_service.GetMinioService()
	if err := m.LoadConfig(configs.GetConfig()); err != nil {
		logg.Fatal("Could not load the MinIO configuration due to an error: %v", err)
	}
	if err := m.Init(); err != nil {
		logg.Fatal("Could not initialize the MinIO service due to an error: %v", err)
	}
	if err := m.AssertBucket(); err != nil {
		logg.Fatal("Could not assert the MinIO bucket due to an error: %v", err)
	}

	// HLS HTTP server
	s := GetHttpServer()
	if err := s.LoadConfig(configs.GetConfig()); err != nil {
		logg.Error("Could not load HLS HTTP server config due to an error: %v", err)
	}
	if err := s.Init().Listen(); err != nil {
		logg.Fatal("Could not start the HTTP server due to an error: %v", err)
	}
}
