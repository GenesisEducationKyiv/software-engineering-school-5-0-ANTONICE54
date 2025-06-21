package providers

import (
	"context"
	"weather-forecast/internal/domain/models"
)

type (
	WeatherProvider interface {
		GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error)
	}

	ChainLinkLink interface {
		WeatherProvider
		SetNext(section ChainLinkLink)
	}

	ChainLink struct {
		provider    WeatherProvider
		nextSection ChainLinkLink
	}
)

func NewChainLink(provider WeatherProvider) *ChainLink {
	return &ChainLink{
		provider:    provider,
		nextSection: nil,
	}
}

func (c *ChainLink) SetNext(section ChainLinkLink) {
	c.nextSection = section
}

func (c *ChainLink) GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error) {

	weather, err := c.provider.GetWeatherByCity(ctx, city)

	if err != nil {
		if c.nextSection != nil {
			return c.nextSection.GetWeatherByCity(ctx, city)
		}

		return nil, err
	}

	return weather, nil

}
