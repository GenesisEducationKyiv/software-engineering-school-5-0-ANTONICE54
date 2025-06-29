package cache

import (
	"context"
	"encoding/json"
	"time"
	"weather-forecast/internal/infrastructure/logger"

	"github.com/redis/go-redis/v9"
)

const RedisTimeout = 5 * time.Second

type (
	MetricsRecorder interface {
		RecordCacheHit()
		RecordCacheMiss()
		RecordCacheError()
	}
	Redis struct {
		client  *redis.Client
		metrics MetricsRecorder
		logger  logger.Logger
	}
)

func NewRedis(url string, metrics MetricsRecorder, logger logger.Logger) *Redis {
	opt, err := redis.ParseURL(url)
	if err != nil {
		logger.Fatalf("parse Redis URL: %s", err.Error())
	}
	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), RedisTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Fatalf("connect to Redis: %s", err.Error())

	}

	return &Redis{
		client:  client,
		metrics: metrics,
		logger:  logger,
	}
}

func (c *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		c.logger.Warnf("Marshal cache:%s", err.Error())
		return err
	}

	if err := c.client.Set(ctx, key, data, expiration).Err(); err != nil {
		c.metrics.RecordCacheError()
		c.logger.Warnf("Set cache key %s:%s", key, err.Error())
		return err
	}

	return nil
}
func (c *Redis) Get(ctx context.Context, key string, value interface{}) error {
	res, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			c.logger.Infof("Cache miss for key %s", key)
			c.metrics.RecordCacheMiss()
			return err
		}

		c.metrics.RecordCacheError()
		c.logger.Warnf("Get cache key %s:%s", key, err.Error())
		return err
	}

	if err := json.Unmarshal([]byte(res), value); err != nil {
		c.logger.Warnf("Unmarshal cache:%s", err.Error())
		return err
	}
	c.metrics.RecordCacheHit()

	return nil
}
