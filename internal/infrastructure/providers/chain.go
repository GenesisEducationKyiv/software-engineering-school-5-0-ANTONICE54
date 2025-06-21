package providers

import (
	"context"
	"weather-forecast/internal/domain/models"
)

type (
	WeatherProvider interface {
		GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error)
	}

	ChainSection interface {
		WeatherProvider
		SetNext(section ChainSection)
	}

	WeatherChain struct {
		provider    WeatherProvider
		nextSection ChainSection
	}
)

func NewWeatherChain(provider WeatherProvider) *WeatherChain {
	return &WeatherChain{
		provider:    provider,
		nextSection: nil,
	}
}

func (c *WeatherChain) SetNext(section ChainSection) {
	c.nextSection = section
}

func (c *WeatherChain) GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error) {

	weather, err := c.provider.GetWeatherByCity(ctx, city)

	if err != nil {
		if c.nextSection != nil {
			return c.nextSection.GetWeatherByCity(ctx, city)
		}

		return nil, err
	}

	return weather, nil

}
