package main

import (
	"os"
	"os/signal"

	"github.com/cbstorm/wyrstream/lib/logger"
)

func main() {
	logg := logger.NewLogger("AUTH_SERVICE")
	logg.Info("Auth service started.")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logg.Info("Control service shutting down...")
	os.Exit(0)
}
