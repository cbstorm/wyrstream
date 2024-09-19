package main

import (
	"os"
	"strconv"
	"sync"

	"github.com/cbstorm/wyrstream/lib/logger"
)

var config *Config
var config_sync sync.Once

func GetConfig() *Config {
	if config == nil {
		config_sync.Do(func() {
			config = &Config{
				logger: logger.NewLogger("CONTROL_SVC_CONFIG"),
			}
		})
	}
	return config

}

type Config struct {
	logger         *logger.Logger
	HLS_PUBLIC_URL string
	HLS_HTTP_HOST  string
	HLS_HTTP_PORT  uint16
}

func (c *Config) Load() error {
	//Http server
	c.HLS_HTTP_HOST = os.Getenv("HLS_HTTP_HOST")
	port, err := strconv.ParseUint(os.Getenv("HLS_HTTP_PORT"), 10, 16)
	if err != nil {
		return err
	}
	c.HLS_HTTP_PORT = uint16(port)
	c.HLS_PUBLIC_URL = os.Getenv("HLS_PUBLIC_URL")
	return nil
}
