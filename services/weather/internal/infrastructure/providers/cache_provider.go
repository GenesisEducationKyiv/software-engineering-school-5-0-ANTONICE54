package providers

import (
	"context"
	"weather-forecast/pkg/apperrors"
	"weather-forecast/pkg/logger"
	"weather-service/internal/domain/models"
	"weather-service/internal/infrastructure/errors"
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

	cachedWeather := &models.Weather{}
	err := p.cache.Get(ctx, city, cachedWeather)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			switch appErr.Code.String() {
			case errors.CacheMissError.Code.String():
				p.metrics.RecordCacheMiss()
			case errors.CacheError.Code.String():
				p.metrics.RecordCacheError()
			}
		}
		return nil, err
	}
	p.metrics.RecordCacheHit()
	return cachedWeather, nil
}
