package providers

import (
	"context"
	"time"
	"weather-forecast/pkg/logger"
	"weather-service/internal/domain/models"
)

const cacheTTL = 10 * time.Minute

type (
	CacheWriter interface {
		Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	}

	CacheErrorRecorder interface {
		RecordCacheError()
	}

	CacheDecorator struct {
		provider WeatherProvider
		cache    CacheWriter
		metrics  CacheErrorRecorder
		logger   logger.Logger
	}
)

func NewCacheDecorator(provider WeatherProvider, cache CacheWriter, metrics CacheErrorRecorder, logger logger.Logger) *CacheDecorator {
	return &CacheDecorator{
		provider: provider,
		cache:    cache,
		metrics:  metrics,
		logger:   logger,
	}
}

func (d *CacheDecorator) GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error) {

	weather, err := d.provider.GetWeatherByCity(ctx, city)

	if err != nil {
		return nil, err
	}

	if err := d.cache.Set(ctx, city, weather, cacheTTL); err != nil {
		d.metrics.RecordCacheError()
		d.logger.Warnf("Cache weather %s:", err.Error())
	}

	return weather, nil

}
