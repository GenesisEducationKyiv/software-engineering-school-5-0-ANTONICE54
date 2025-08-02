package providers

import (
	"context"
	"time"
	"weather-forecast/pkg/logger"
	"weather-service/internal/domain/models"
	"weather-service/internal/domain/usecases"
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
		provider usecases.WeatherProvider
		cache    CacheWriter
		metrics  CacheErrorRecorder
		logger   logger.Logger
	}
)

func NewCacheDecorator(provider usecases.WeatherProvider, cache CacheWriter, metrics CacheErrorRecorder, logger logger.Logger) *CacheDecorator {
	return &CacheDecorator{
		provider: provider,
		cache:    cache,
		metrics:  metrics,
		logger:   logger,
	}
}

func (d *CacheDecorator) GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error) {

	log := d.logger.WithContext(ctx)

	log.Debugf("Getting weather and caching for city: %s", city)

	weather, err := d.provider.GetWeatherByCity(ctx, city)

	if err != nil {
		return nil, err
	}

	if err := d.cache.Set(ctx, city, weather, cacheTTL); err != nil {
		d.metrics.RecordCacheError()
		log.Warnf("Failed to cache weather for city %s: %v", city, err)
	}

	log.Debugf("Weather cached successfully for city: %s", city)

	return weather, nil

}
