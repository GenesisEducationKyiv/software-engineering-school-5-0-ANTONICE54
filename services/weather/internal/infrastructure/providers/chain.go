package providers

import (
	"context"
	"weather-service/internal/domain/models"
)

type (
	WeatherProvider interface {
		GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error)
	}

	WeatherChainLink interface {
		WeatherProvider
		SetNext(section WeatherChainLink)
	}

	WeatherLink struct {
		provider    WeatherProvider
		nextSection WeatherChainLink
	}
)

func NewWeatherLink(provider WeatherProvider) *WeatherLink {
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
