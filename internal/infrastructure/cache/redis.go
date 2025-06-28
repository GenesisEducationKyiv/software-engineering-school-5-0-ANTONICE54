package cache

import (
	"context"
	"time"
	"weather-forecast/internal/infrastructure/logger"

	"github.com/redis/go-redis/v9"
)

const RedisTimeout = 5 * time.Second

type (
	Redis struct {
		*redis.Client
		logger logger.Logger
	}
)

func NewRedis(url string, logger logger.Logger) *Redis {
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
		Client: client,
		logger: logger,
	}
}

func (c *Redis) Set(ctx context.Context, key string) {}
func (c *Redis) Get(ctx context.Context, key string) {}
