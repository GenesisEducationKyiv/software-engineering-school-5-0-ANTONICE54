package services

import (
	"context"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/logger"
)

type (
	WeatherProvider interface {
		GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error)
	}
	WeatherService struct {
		weatherProvider WeatherProvider
		logger          logger.Logger
	}
)

func NewWeatherService(weatherProvider WeatherProvider, logger logger.Logger) *WeatherService {
	return &WeatherService{
		weatherProvider: weatherProvider,
		logger:          logger,
	}
}

func (s *WeatherService) GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error) {
	weather, err := s.weatherProvider.GetWeatherByCity(ctx, city)
	if err != nil {
		return nil, err
	}
	return weather, nil
}
