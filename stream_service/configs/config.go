package configs

import (
	"os"
	"sync"

	"github.com/cbstorm/wyrstream/lib/logger"
)

var instance *Config
var instance_sync sync.Once

func GetConfig() *Config {
	if instance == nil {
		instance_sync.Do(func() {
			instance = &Config{
				logger: logger.NewLogger("STREAM_SVC_CONFIG"),
			}
		})
	}
	return instance

}

type Config struct {
	logger     *logger.Logger
	ADDR       string
	PUBLIC_URL string
}

func (c *Config) Load() error {
	c.ADDR = os.Getenv("ADDR")
	c.PUBLIC_URL = os.Getenv("PUBLIC_URL")
	return nil
}
