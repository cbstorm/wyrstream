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
	logger         *logger.Logger
	APP_NAME       string
	JWT_SECRET_KEY string

	MONGODB_URL     string
	MONGODB_DB_NAME string

	NATS_CORE_USERNAME    string
	NATS_CORE_PASSWORD    string
	NATS_CORE_HOST        string
	NATS_CORE_PORT        uint16
	NATS_CORE_QUEUE_GROUP string

	REDIS_USERNAME   string
	REDIS_PASSWORD   string
	REDIS_HOST       string
	REDIS_PORT       uint16
	REDIS_KEY_PREFIX string
}

func (c *Config) Load() error {
	c.APP_NAME = os.Getenv("APP_NAME")
	c.JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
	// MongoDB
	c.MONGODB_URL = os.Getenv("MONGODB_URL")
	c.MONGODB_DB_NAME = os.Getenv("MONGODB_DB_NAME")
	// NATS
	c.NATS_CORE_USERNAME = os.Getenv("NATS_CORE_USERNAME")
	c.NATS_CORE_PASSWORD = os.Getenv("NATS_CORE_PASSWORD")
	c.NATS_CORE_HOST = os.Getenv("NATS_CORE_HOST")
	nats_core_port, err := strconv.ParseUint(os.Getenv("NATS_CORE_PORT"), 10, 32)
	if err != nil {
		return err
	}
	c.NATS_CORE_PORT = uint16(nats_core_port)
	c.NATS_CORE_QUEUE_GROUP = os.Getenv("NATS_CORE_QUEUE_GROUP")
	// Redis
	c.REDIS_USERNAME = os.Getenv("REDIS_USERNAME")
	c.REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
	c.REDIS_HOST = os.Getenv("REDIS_HOST")
	redis_port, err := strconv.ParseUint(os.Getenv("REDIS_PORT"), 10, 16)
	if err != nil {
		return err
	}
	c.REDIS_PORT = uint16(redis_port)
	c.REDIS_KEY_PREFIX = os.Getenv("REDIS_KEY_PREFIX")
	return nil
}
