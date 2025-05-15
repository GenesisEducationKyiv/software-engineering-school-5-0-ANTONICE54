package services

import (
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/logger"
)

type (
	WeatherProviderI interface {
		GetWeatherByCity(city string) (*models.Weather, error)
	}
	WeatherService struct {
		weatherProvider WeatherProviderI
		logger          logger.Logger
	}
)

func NewWeatherService(weatherProvider WeatherProviderI, logger logger.Logger) *WeatherService {
	return &WeatherService{
		weatherProvider: weatherProvider,
		logger:          logger,
	}
}

func (s *WeatherService) GetWeatherByCity(city string) (*models.Weather, error) {
	weather, err := s.weatherProvider.GetWeatherByCity(city)
	if err != nil {
		return nil, err
	}
	return weather, nil
}
