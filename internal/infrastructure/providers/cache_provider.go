package providers

import (
	"context"
	"time"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/logger"
)

const cacheTTL = 10 * time.Minute

type (
	Cacher interface {
		Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
		Get(ctx context.Context, key string, value interface{}) error
	}

	CacheWeatherProvider struct {
		cache    Cacher
		provider WeatherProvider
		logger   logger.Logger
	}
)

func NewCacheWeather(cache Cacher, provider WeatherProvider, logger logger.Logger) *CacheWeatherProvider {
	return &CacheWeatherProvider{
		cache:    cache,
		provider: provider,
		logger:   logger,
	}
}

func (p *CacheWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error) {

	cachedWeather := &models.Weather{}
	err := p.cache.Get(ctx, city, cachedWeather)
	if err == nil {
		return cachedWeather, nil
	}

	weather, err := p.provider.GetWeatherByCity(ctx, city)
	if err != nil {
		return nil, err
	}

	if err := p.cache.Set(ctx, city, weather, cacheTTL); err != nil {
		p.logger.Warnf("Cache weather %s:", err.Error())
	}

	return weather, nil
}
