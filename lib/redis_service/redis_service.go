package redis_service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/redis/go-redis/v9"
)

type IRedisConfig interface {
	LoadRedisConfig() error
	REDIS_USERNAME() string
	REDIS_PASSWORD() string
	REDIS_HOST() string
	REDIS_PORT() uint16
	REDIS_KEY_PREFIX() string
}

var redis_service *RedisService
var redis_service_sync sync.Once

func GetRedisService() *RedisService {
	if redis_service == nil {
		redis_service_sync.Do(func() {
			redis_service = &RedisService{
				logger: logger.NewLogger("REDIS_SERVICE"),
			}
		})
	}
	return redis_service
}

type RedisService struct {
	rdb        *redis.Client
	key_prefix string
	logger     *logger.Logger
	config     IRedisConfig
}

func (i *RedisService) LoadConfig(config IRedisConfig) error {
	if err := config.LoadRedisConfig(); err != nil {
		return err
	}
	i.config = config
	return nil
}

func (i *RedisService) Connect() error {
	rdb := redis.NewClient(&redis.Options{
		Username: i.config.REDIS_USERNAME(),
		Password: i.config.REDIS_PASSWORD(),
		Addr:     fmt.Sprintf("%s:%d", i.config.REDIS_HOST(), i.config.REDIS_PORT()),
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	i.rdb = rdb
	i.key_prefix = i.config.REDIS_KEY_PREFIX()
	i.logger.Info("Connected to redis successfully")
	return nil
}

func (i *RedisService) Close() error {
	return i.rdb.Close()
}

func (i *RedisService) Set(key string, value interface{}) error {
	byteValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, err = i.rdb.Set(ctx, i.getKey(key), string(byteValue), 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func (i *RedisService) SetWithTtl(key string, value interface{}, second time.Duration) error {
	byteValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, err = i.rdb.Set(ctx, i.getKey(key), string(byteValue), second).Result()
	if err != nil {
		return err
	}
	return nil
}

func (i *RedisService) Get(key string, out interface{}) (error, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	value, err := i.rdb.Get(ctx, i.getKey(key)).Result()
	if err != nil && err == redis.Nil {
		return nil, true
	}
	if err != nil {
		return err, false
	}

	err = json.Unmarshal([]byte(value), out)
	if err != nil {
		return err, false
	}
	return nil, false
}

func (i *RedisService) getKey(k string) string {
	return fmt.Sprintf("%s:%s", i.key_prefix, k)
}
