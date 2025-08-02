package usecases

import (
	"context"
	"weather-forecast/pkg/logger"
	"weather-service/internal/domain/models"
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
	log := s.logger.WithContext(ctx)

	log.Infof("Getting weather for city: %s", city)

	weather, err := s.weatherProvider.GetWeatherByCity(ctx, city)
	if err != nil {
		log.Errorf("Failed to get weather for city %s: %v", city, err)

		return nil, err
	}

	log.Infof("Weather retrieved successfully for city: %s", city)

	return weather, nil
}
