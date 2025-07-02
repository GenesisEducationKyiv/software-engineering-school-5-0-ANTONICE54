package usecases

import (
	"context"
	"weather-forecast/internal/domain/models"
)

type WeatherProvider interface {
	GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error)
}
