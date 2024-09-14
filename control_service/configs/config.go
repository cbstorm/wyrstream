package configs

import (
	"os"
	"strconv"
	"sync"

	"github.com/cbstorm/wyrstream/lib/logger"
)

var instance *Config
var instance_sync sync.Once

func GetConfig() *Config {
	if instance == nil {
		instance_sync.Do(func() {
			instance = &Config{
				logger: logger.NewLogger("CONTROL_SVC_CONFIG"),
			}
		})
	}
	return instance

}

type Config struct {
	logger    *logger.Logger
	HTTP_HOST string
	HTTP_PORT uint16
}

func (c *Config) Load() error {
	//Http server
	c.HTTP_HOST = os.Getenv("HTTP_HOST")
	port, err := strconv.ParseUint(os.Getenv("HTTP_PORT"), 10, 16)
	if err != nil {
		return err
	}
	c.HTTP_PORT = uint16(port)
	return nil
}
