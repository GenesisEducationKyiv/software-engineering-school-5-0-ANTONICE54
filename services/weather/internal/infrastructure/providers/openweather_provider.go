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

	log := p.logger.WithContext(ctx)

	log.Debugf("Requesting weather data from OpenWeather for city: %s", city)

	weatherResponse, err := p.client.GetWeather(ctx, city)
	log.Debugf("Processing OpenWeather response for city: %s", city)
	if err != nil {
		return nil, err
	}

	weatherDesc := ""

	if len(weatherResponse.Weather) > 0 {
		weatherDesc = weatherResponse.Weather[0].Description
	} else {
		log.Warnf("OpenWeather did not provide weather description for city: %s", city)
	}

	result := models.Weather{
		Temperature: weatherResponse.Main.Temperature,
		Humidity:    weatherResponse.Main.Humidity,
		Description: weatherDesc,
	}

	log.Infof("OpenWeather data processed successfully for city: %s", city)

	return &result, nil

}
