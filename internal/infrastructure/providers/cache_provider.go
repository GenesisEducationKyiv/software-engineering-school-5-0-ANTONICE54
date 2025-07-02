package providers

import (
	"context"
	"time"
	"weather-forecast/internal/domain/models"
	infraerrors "weather-forecast/internal/infrastructure/errors"
	"weather-forecast/internal/infrastructure/logger"
	"weather-forecast/pkg/apperrors"
)

type (
	MetricsRecorder interface {
		RecordCacheHit()
		RecordCacheMiss()
		RecordCacheError()
	}

	Cacher interface {
		Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
		Get(ctx context.Context, key string, value interface{}) error
	}

	CacheWeatherProvider struct {
		cache   Cacher
		metrics MetricsRecorder
		logger  logger.Logger
	}
)

func NewCacheWeather(cache Cacher, metrics MetricsRecorder, logger logger.Logger) *CacheWeatherProvider {
	return &CacheWeatherProvider{
		cache:   cache,
		metrics: metrics,
		logger:  logger,
	}
}

func (p *CacheWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error) {

	cachedWeather := &models.Weather{}
	err := p.cache.Get(ctx, city, cachedWeather)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			switch appErr.Code.String() {
			case infraerrors.CacheMissError.Code.String():
				p.metrics.RecordCacheMiss()
			case infraerrors.CacheError.Code.String():
				p.metrics.RecordCacheError()
			}
		}
		return nil, err
	}
	p.metrics.RecordCacheHit()
	return cachedWeather, nil
}
