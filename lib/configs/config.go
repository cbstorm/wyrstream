package configs

import (
	"os"
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
	logger          *logger.Logger
	MONGODB_URL     string
	MONGODB_DB_NAME string
}

func (c *Config) Load() error {
	c.MONGODB_URL = os.Getenv("MONGODB_URL")
	c.MONGODB_DB_NAME = os.Getenv("MONGODB_DB_NAME")
	return nil
}
