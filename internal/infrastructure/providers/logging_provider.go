package providers

import (
	"context"
	"encoding/json"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/logger"
)

type (
	DescriptiveWeatherProvider interface {
		WeatherProvider
		Name() string
	}

	LoggingProvider struct {
		wrapped DescriptiveWeatherProvider
		logger  logger.Logger
	}
)

func NewLogging(provider DescriptiveWeatherProvider, logger logger.Logger) *LoggingProvider {
	return &LoggingProvider{
		wrapped: provider,
		logger:  logger,
	}
}

func (p *LoggingProvider) GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error) {

	weather, err := p.wrapped.GetWeatherByCity(ctx, city)

	if err != nil {
		//TODO: add logging logic

		return nil, err

	}

	weatherJSON, err := json.Marshal(weather)
	if err != nil {
		p.logger.Warnf("%s - failed to marshal weather data for city '%s':%s", p.wrapped.Name(), city, err.Error())
		return weather, nil

	}
	p.logger.Infof("%s - got weather for city '%s':%s", p.wrapped.Name(), city, string(weatherJSON))

	return weather, nil

}
