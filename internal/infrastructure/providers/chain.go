package providers

import (
	"context"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/services"
)

type (
	ChainSection interface {
		services.WeatherProvider
		SetNext(section ChainSection)
	}

	WeatherProviderChain struct {
		provider    services.WeatherProvider
		nextSection ChainSection
	}
)

func NewWeatherProviderChain(provider services.WeatherProvider) *WeatherProviderChain {
	return &WeatherProviderChain{
		provider:    provider,
		nextSection: nil,
	}
}

func (c *WeatherProviderChain) SetNext(section ChainSection) {
	c.nextSection = section
}

func (c *WeatherProviderChain) GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error) {

	weather, err := c.provider.GetWeatherByCity(ctx, city)

	if err != nil {
		if c.nextSection != nil {
			return c.nextSection.GetWeatherByCity(ctx, city)
		}

		return nil, err
	}

	return weather, nil

}
