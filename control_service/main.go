package main

import (
	local_configs "github.com/cbstorm/wyrstream/control_service/configs"
	"github.com/cbstorm/wyrstream/control_service/http_server"
	"github.com/cbstorm/wyrstream/lib/configs"
	"github.com/cbstorm/wyrstream/lib/database"
	"github.com/cbstorm/wyrstream/lib/logger"
	natsservice "github.com/cbstorm/wyrstream/lib/nats_service"
	"github.com/cbstorm/wyrstream/lib/redis_service"
)

func main() {
	logg := logger.NewLogger("CONTROL_SVC")
	if err := local_configs.GetConfig().Load(); err != nil {
		logg.Fatal("Could not load local configs with err: %v", err)
	}
	if err := configs.GetConfig().Load(); err != nil {
		logg.Fatal("Could not load configs with err: %v", err)
	}
	if err := database.GetDatabase().Connect(); err != nil {
		logg.Fatal("Could not connect to database %v", err)
	}
	if err := redis_service.GetRedisService().Connect(); err != nil {
		logg.Fatal("Could not connect to redis server with err: %v", err)
	}
	if err := natsservice.GetNATSService().Connect(); err != nil {
		logg.Fatal("Could not connect to NATS server with err: %v", err)
	}
	if err := http_server.GetHttpServer().AddRoutes().Listen(); err != nil {
		logg.Fatal("Could not start http server with err: %v", err)
	}
}
