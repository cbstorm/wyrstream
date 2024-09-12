package main

import (
	"github.com/cbstorm/wyrstream/lib/configs"
	"github.com/cbstorm/wyrstream/lib/database"
	"github.com/cbstorm/wyrstream/lib/logger"
)

func main() {
	logg := logger.NewLogger("CONTROL_SVC")
	if err := configs.GetConfig().Load(); err != nil {
		logg.Fatal("Could not load config with err: %v", err)
	}
	if err := database.GetDatabase().Connect(); err != nil {
		logg.Fatal("Could not connect to database %v", err)
	}
	logg.Info("[control_service] started")
}
