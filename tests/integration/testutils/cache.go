package testutils

import (
	"context"
	"sync"
	"time"
	"weather-forecast/internal/domain/models"
	infraerrors "weather-forecast/internal/infrastructure/errors"
)

type (
	cacheItem struct {
		value      interface{}
		expiration time.Time
	}

	InMemoryCache struct {
		data map[string]cacheItem
		mu   *sync.Mutex
	}
)

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]cacheItem),
		mu:   &sync.Mutex{},
	}
}

func (c *InMemoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(expiration),
	}

	return nil
}

func (c *InMemoryCache) Get(ctx context.Context, key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.data[key]
	if !exists {
		return infraerrors.CacheMissError
	}

	if time.Now().After(item.expiration) {
		delete(c.data, key)
		return infraerrors.CacheMissError
	}

	if weather, ok := value.(*models.Weather); ok {
		if cachedWeather, ok := item.value.(*models.Weather); ok {
			*weather = *cachedWeather
			return nil
		}
	}

	return infraerrors.CacheError
}

func (c *InMemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]cacheItem)
}

func (c *InMemoryCache) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.data)
}
