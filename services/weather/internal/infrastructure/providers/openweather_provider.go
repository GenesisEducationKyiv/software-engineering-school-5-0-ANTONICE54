package providers

import (
	"context"
	"weather-service/internal/domain/models"
	"weather-service/internal/infrastructure/clients/openweather"

	"weather-forecast/pkg/logger"
)

type (
	OpenWeatherProvider struct {
		client *openweather.OpenWeatherClient
		logger logger.Logger
	}
)

func NewOpenWeatherProvider(client *openweather.OpenWeatherClient, logger logger.Logger) *OpenWeatherProvider {
	return &OpenWeatherProvider{
		client: client,
		logger: logger,
	}
}

func (p *OpenWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error) {

	weatherResponse, err := p.client.GetWeather(ctx, city)
	if err != nil {
		return nil, err
	}

	weatherDesc := ""

	if len(weatherResponse.Weather) > 0 {
		weatherDesc = weatherResponse.Weather[0].Description
	} else {
		p.logger.Warnf("OpenWeather did not provide weather description for city: %s", city)
	}

	result := models.Weather{
		Temperature: weatherResponse.Main.Temperature,
		Humidity:    weatherResponse.Main.Humidity,
		Description: weatherDesc,
	}

	return &result, nil

}
