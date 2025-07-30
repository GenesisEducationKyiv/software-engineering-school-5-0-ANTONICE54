package providers

import (
	"context"
	"weather-service/internal/domain/models"
	"weather-service/internal/infrastructure/clients/weatherapi"

	"weather-forecast/pkg/logger"
)

type (
	WeatherAPIProvider struct {
		client *weatherapi.WeatherAPIClient
		logger logger.Logger
	}
)

func NewWeatherAPIProvider(client *weatherapi.WeatherAPIClient, logger logger.Logger) *WeatherAPIProvider {
	return &WeatherAPIProvider{

		client: client,
		logger: logger,
	}
}

func (p *WeatherAPIProvider) GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error) {

	weatherResponse, err := p.client.GetWeather(ctx, city)

	if err != nil {
		return nil, err
	}

	result := models.Weather{
		Temperature: weatherResponse.Current.TempC,
		Humidity:    weatherResponse.Current.Humidity,
		Description: weatherResponse.Current.Condition.Text,
	}

	return &result, nil
}
