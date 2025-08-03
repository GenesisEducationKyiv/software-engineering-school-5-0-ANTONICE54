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

	log := p.logger.WithContext(ctx)

	log.Debugf("Requesting weather data from WeatherAPI for city: %s", city)

	weatherResponse, err := p.client.GetWeather(ctx, city)
	log.Debugf("Processing WeatherAPU response for city: %s", city)

	if err != nil {
		return nil, err
	}

	result := models.Weather{
		Temperature: weatherResponse.Current.TempC,
		Humidity:    weatherResponse.Current.Humidity,
		Description: weatherResponse.Current.Condition.Text,
	}

	log.Infof("WeatherAPI data processed successfully for city: %s", city)

	return &result, nil
}
