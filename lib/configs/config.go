package configs

import (
	"os"
	"strconv"
	"sync"

	"github.com/cbstorm/wyrstream/lib/logger"
)

var cfg *Config
var cfgSync sync.Once

func GetConfig() *Config {
	if cfg == nil {
		cfgSync.Do(func() {
			cfg = &Config{
				logger: logger.NewLogger("CONFIG"),
			}
		})
	}
	return cfg

}

type Config struct {
	logger                *logger.Logger
	MONGODB_URL           string
	MONGODB_DB_NAME       string
	NATS_CORE_USERNAME    string
	NATS_CORE_PASSWORD    string
	NATS_CORE_HOST        string
	NATS_CORE_PORT        uint
	NATS_CORE_QUEUE_GROUP string
}

func (c *Config) Load() error {
	c.MONGODB_URL = os.Getenv("MONGODB_URL")
	c.MONGODB_DB_NAME = os.Getenv("MONGODB_DB_NAME")
	c.NATS_CORE_USERNAME = os.Getenv("NATS_CORE_USERNAME")
	c.NATS_CORE_PASSWORD = os.Getenv("NATS_CORE_PASSWORD")
	c.NATS_CORE_HOST = os.Getenv("NATS_CORE_HOST")
	nats_core_port, err := strconv.ParseUint(os.Getenv("NATS_CORE_PORT"), 10, 32)
	if err != nil {
		return err
	}
	c.NATS_CORE_PORT = uint(nats_core_port)
	c.NATS_CORE_QUEUE_GROUP = os.Getenv("NATS_CORE_QUEUE_GROUP")
	return nil
}
