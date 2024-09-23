package configs

import (
	"errors"
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

	// Http server control_service
	http_host string
	http_port uint16

	// Database
	mongodb_url  string
	mongodb_name string

	// NATS
	nats_core_username    string
	nats_core_password    string
	nats_core_host        string
	nats_core_port        uint16
	nats_core_queue_group string

	// Redis
	redis_username   string
	redis_password   string
	redis_host       string
	redis_port       uint16
	redis_key_prefix string

	// HLS
	hls_public_url string
	hls_http_host  string
	hls_http_port  uint16

	// Stream server
	stream_server_addr       string
	stream_server_public_url string

	// MinIO
	minio_host        string
	minio_port        uint16
	minio_access_key  string
	minio_secret_key  string
	minio_bucket_name string
	minio_public_url  string
}

func (c *Config) Load() error {
	c.APP_NAME = os.Getenv("APP_NAME")
	if c.APP_NAME == "" {
		return errors.New("required APP_NAME")
	}
	c.JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
	if c.JWT_SECRET_KEY == "" {
		return errors.New("required JWT_SECRET_KEY")
	}
	return nil
}

// HTTP server control_service
func (c *Config) LoadHttpConfig() error {
	c.http_host = os.Getenv("HTTP_HOST")
	if c.http_host == "" {
		return errors.New("required HTTP_HOST")
	}
	port, err := strconv.ParseUint(os.Getenv("HTTP_PORT"), 10, 16)
	if err != nil {
		return err
	}
	c.http_port = uint16(port)
	return nil
}
func (c *Config) HTTP_HOST() string {
	return c.http_host
}
func (c *Config) HTTP_PORT() uint16 {
	return c.http_port
}

// Database config
func (c *Config) LoadDatabaseConfig() error {
	c.mongodb_url = os.Getenv("MONGODB_URL")
	if c.mongodb_url == "" {
		return errors.New("required MONGODB_URL")
	}
	c.mongodb_name = os.Getenv("MONGODB_DB_NAME")
	if c.mongodb_name == "" {
		return errors.New("required MONGODB_DB_NAME")
	}
	return nil
}

func (c *Config) MONGODB_URL() string {
	return c.mongodb_url
}

func (c *Config) MONGODB_DB_NAME() string {
	return c.mongodb_name
}

// Redis config
func (c *Config) LoadRedisConfig() error {
	c.redis_username = os.Getenv("REDIS_USERNAME")
	if c.redis_username == "" {
		return errors.New("required REDIS_USERNAME")
	}
	c.redis_password = os.Getenv("REDIS_PASSWORD")
	if c.redis_password == "" {
		return errors.New("required REDIS_PASSWORD")
	}
	c.redis_host = os.Getenv("REDIS_HOST")
	if c.redis_host == "" {
		return errors.New("required REDIS_HOST")
	}
	redis_port, err := strconv.ParseUint(os.Getenv("REDIS_PORT"), 10, 16)
	if err != nil {
		return err
	}
	c.redis_port = uint16(redis_port)
	c.redis_key_prefix = os.Getenv("REDIS_KEY_PREFIX")
	return nil
}
func (c *Config) REDIS_USERNAME() string {
	return c.redis_username
}
func (c *Config) REDIS_PASSWORD() string {
	return c.redis_password
}
func (c *Config) REDIS_HOST() string {
	return c.redis_host
}
func (c *Config) REDIS_PORT() uint16 {
	return c.redis_port
}
func (c *Config) REDIS_KEY_PREFIX() string {
	return c.redis_key_prefix
}

// NATS Config
func (c *Config) LoadNATSCoreConfig() error {
	c.nats_core_username = os.Getenv("NATS_CORE_USERNAME")
	if c.nats_core_username == "" {
		return errors.New("required NATS_CORE_USERNAME")
	}
	c.nats_core_password = os.Getenv("NATS_CORE_PASSWORD")
	if c.nats_core_password == "" {
		return errors.New("required NATS_CORE_PASSWORD")
	}
	c.nats_core_host = os.Getenv("NATS_CORE_HOST")
	if c.nats_core_host == "" {
		return errors.New("required NATS_CORE_HOST")
	}
	nats_core_port, err := strconv.ParseUint(os.Getenv("NATS_CORE_PORT"), 10, 32)
	if err != nil {
		return err
	}
	c.nats_core_port = uint16(nats_core_port)
	c.nats_core_queue_group = os.Getenv("NATS_CORE_QUEUE_GROUP")
	return nil
}

func (c *Config) NATS_CORE_USERNAME() string {
	return c.nats_core_username
}
func (c *Config) NATS_CORE_PASSWORD() string {
	return c.nats_core_password
}
func (c *Config) NATS_CORE_HOST() string {
	return c.nats_core_host
}
func (c *Config) NATS_CORE_PORT() uint16 {
	return c.nats_core_port
}
func (c *Config) NATS_CORE_QUEUE_GROUP() string {
	return c.nats_core_queue_group
}

// Stream server
func (c *Config) LoadStreamServerConfig() error {
	c.stream_server_addr = os.Getenv("STREAM_SERVER_ADDR")
	if c.stream_server_addr == "" {
		return errors.New("required STREAM_SERVER_ADDR")
	}
	c.stream_server_public_url = os.Getenv("STREAM_SERVER_PUBLIC_URL")
	if c.stream_server_public_url == "" {
		return errors.New("required STREAM_SERVER_PUBLIC_URL")
	}
	return nil
}
func (c *Config) SERVER_ADDRESS() string {
	return c.stream_server_addr
}
func (c *Config) SERVER_PUBLIC_URL() string {
	return c.stream_server_public_url
}

// HLS config
func (c *Config) LoadHLSHttpServerConfig() error {
	c.hls_http_host = os.Getenv("HLS_HTTP_HOST")
	if c.hls_http_host == "" {
		return errors.New("required HLS_HTTP_HOST")
	}
	port, err := strconv.ParseUint(os.Getenv("HLS_HTTP_PORT"), 10, 16)
	if err != nil {
		return err
	}
	c.hls_http_port = uint16(port)
	c.hls_public_url = os.Getenv("HLS_PUBLIC_URL")
	if c.hls_public_url == "" {
		return errors.New("required HLS_PUBLIC_URL")
	}
	return nil
}

func (c *Config) HLS_HTTP_HOST() string {
	return c.hls_http_host
}
func (c *Config) HLS_HTTP_PORT() uint16 {
	return c.hls_http_port
}
func (c *Config) HLS_PUBLIC_URL() string {
	return c.hls_public_url
}

// Minio config
func (c *Config) LoadMinioConfig() error {
	c.minio_host = os.Getenv("MINIO_HOST")
	if c.minio_host == "" {
		return errors.New("required MINIO_HOST")
	}
	port, err := strconv.ParseUint(os.Getenv("MINIO_PORT"), 10, 16)
	if err != nil {
		return errors.New("invalid MINIO_PORT")
	}
	c.minio_port = uint16(port)
	c.minio_access_key = os.Getenv("MINIO_ACCESS_KEY")
	if c.minio_access_key == "" {
		return errors.New("required MINIO_ACCESS_KEY")
	}
	c.minio_secret_key = os.Getenv("MINIO_SECRET_KEY")
	if c.minio_secret_key == "" {
		return errors.New("required MINIO_SECRET_KEY")
	}
	c.minio_bucket_name = os.Getenv("MINIO_BUCKET_NAME")
	if c.minio_bucket_name == "" {
		return errors.New("required MINIO_BUCKET_NAME")
	}
	c.minio_public_url = os.Getenv("MINIO_PUBLIC_URL")
	if c.minio_public_url == "" {
		return errors.New("required MINIO_PUBLIC_URL")
	}
	return nil
}
func (c *Config) MINIO_HOST() string {
	return c.minio_host
}
func (c *Config) MINIO_PORT() uint16 {
	return c.minio_port
}
func (c *Config) MINIO_ACCESS_KEY() string {
	return c.minio_access_key
}
func (c *Config) MINIO_SECRET_KEY() string {
	return c.minio_secret_key
}
func (c *Config) MINIO_BUCKET_NAME() string {
	return c.minio_bucket_name
}

func (c *Config) MINIO_PUBLIC_URL() string {
	return c.minio_public_url
}
