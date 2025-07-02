package usecases

import (
	"context"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/logger"
)

type (
	WeatherUseCase struct {
		weatherProvider WeatherProvider
		logger          logger.Logger
	}
)

func NewWeatherService(weatherProvider WeatherProvider, logger logger.Logger) *WeatherUseCase {
	return &WeatherUseCase{
		weatherProvider: weatherProvider,
		logger:          logger,
	}
}

func (uc *WeatherUseCase) GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error) {
	weather, err := uc.weatherProvider.GetWeatherByCity(ctx, city)
	if err != nil {
		return nil, err
	}
	return weather, nil
}
