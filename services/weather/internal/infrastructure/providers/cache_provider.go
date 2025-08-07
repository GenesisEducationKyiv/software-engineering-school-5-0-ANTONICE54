package providers

import (
	"context"
	"errors"
	"weather-forecast/pkg/logger"
	"weather-service/internal/domain/models"
	infraerrors "weather-service/internal/infrastructure/errors"
)

type (
	MetricsRecorder interface {
		RecordCacheHit()
		RecordCacheMiss()
		RecordCacheError()
	}

	CacheReader interface {
		Get(ctx context.Context, key string, value interface{}) error
	}

	CacheWeatherProvider struct {
		cache   CacheReader
		metrics MetricsRecorder
		logger  logger.Logger
	}
)

func NewCacheWeather(cache CacheReader, metrics MetricsRecorder, logger logger.Logger) *CacheWeatherProvider {
	return &CacheWeatherProvider{
		cache:   cache,
		metrics: metrics,
		logger:  logger,
	}
}

func (p *CacheWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error) {
	log := p.logger.WithContext(ctx)

	cachedWeather := &models.Weather{}
	err := p.cache.Get(ctx, city, cachedWeather)
	if err != nil {

		if errors.Is(err, infraerrors.ErrCache) {
			log.Errorf("Cache error for city %s: %v", city, err)
			p.metrics.RecordCacheError()
		}

		if errors.Is(err, infraerrors.ErrCacheMiss) {
			p.metrics.RecordCacheMiss()
		}

		return nil, err

	}

	p.metrics.RecordCacheHit()
	return cachedWeather, nil
}
