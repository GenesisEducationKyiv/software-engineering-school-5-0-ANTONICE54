package cache

import (
	"context"
	"encoding/json"
	"time"
	"weather-forecast/pkg/logger"
	"weather-service/internal/infrastructure/errors"

	"github.com/redis/go-redis/v9"
)

const RedisTimeout = 5 * time.Second

type (
	Redis struct {
		client *redis.Client
		logger logger.Logger
	}
)

func NewRedis(url string, logger logger.Logger) (*Redis, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), RedisTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Redis{
		client: client,
		logger: logger,
	}, nil
}

func (c *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		c.logger.Warnf("Marshal cache:%s", err.Error())
		return errors.InternalServerError
	}

	if err := c.client.Set(ctx, key, data, expiration).Err(); err != nil {
		c.logger.Warnf("Set cache key %s:%s", key, err.Error())
		return errors.CacheError
	}

	return nil
}
func (c *Redis) Get(ctx context.Context, key string, value interface{}) error {
	res, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			c.logger.Infof("Cache miss for key %s", key)
			return errors.CacheMissError
		}

		c.logger.Warnf("Get cache key %s:%s", key, err.Error())
		return errors.CacheError
	}

	if err := json.Unmarshal([]byte(res), value); err != nil {
		c.logger.Warnf("Unmarshal cache:%s", err.Error())
		return errors.InternalServerError
	}

	return nil
}
