package providers

import (
	"context"
	"weather-service/internal/domain/models"
	"weather-service/internal/domain/usecases"
)

type (
	WeatherChainLink interface {
		usecases.WeatherProvider
		SetNext(section WeatherChainLink)
	}

	WeatherLink struct {
		provider    usecases.WeatherProvider
		nextSection WeatherChainLink
	}
)

func NewWeatherLink(provider usecases.WeatherProvider) *WeatherLink {
	return &WeatherLink{
		provider:    provider,
		nextSection: nil,
	}
}

func (c *WeatherLink) SetNext(section WeatherChainLink) {
	c.nextSection = section
}

func (c *WeatherLink) GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error) {

	weather, err := c.provider.GetWeatherByCity(ctx, city)

	if err != nil {
		if c.nextSection != nil {
			return c.nextSection.GetWeatherByCity(ctx, city)
		}

		return nil, err
	}

	return weather, nil

}
