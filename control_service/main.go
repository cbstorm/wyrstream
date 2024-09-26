package main

import (
	"github.com/cbstorm/wyrstream/control_service/http_server"
	"github.com/cbstorm/wyrstream/lib/configs"
	"github.com/cbstorm/wyrstream/lib/database"
	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/cbstorm/wyrstream/lib/minio_service"
	"github.com/cbstorm/wyrstream/lib/nats_service"
	"github.com/cbstorm/wyrstream/lib/redis_service"
)

func main() {
	logg := logger.NewLogger("CONTROL_SVC")
	if err := configs.GetConfig().Load(); err != nil {
		logg.Fatal("Could not load config with error: %v", err)
	}
	// DB
	db := database.GetDatabase()
	if err := db.LoadConfig(configs.GetConfig()); err != nil {
		logg.Fatal("Could not load database config with err: %v", err)
	}
	if err := db.Connect(); err != nil {
		logg.Fatal("Could not connect to database %v", err)
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
		logg.Fatal("Could not NATS core config with err: %v", err)
	}
	if err := n.Connect(); err != nil {
		logg.Fatal("Could not connect to NATS server with err: %v", err)
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

	// HTTP server
	s := http_server.GetHttpServer()
	if err := s.LoadConfig(configs.GetConfig()); err != nil {
		logg.Fatal("Could not load http server config with err: %v", err)
	}
	if err := s.Init().AddRoutes().Listen(); err != nil {
		logg.Fatal("Could not start http server with err: %v", err)
	}
}
