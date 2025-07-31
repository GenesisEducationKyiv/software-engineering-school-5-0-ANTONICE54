package clients

import (
	"context"
	"weather-broadcast-service/internal/dto"
	"weather-broadcast-service/internal/mappers"

	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/proto/weather"
)

type (
	WeatherGRPCClient struct {
		weatherGRPC weather.WeatherServiceClient
		logger      logger.Logger
	}
)

func NewWeatherGRPCClient(weatherGRPCClinet weather.WeatherServiceClient, logger logger.Logger) *WeatherGRPCClient {
	return &WeatherGRPCClient{
		weatherGRPC: weatherGRPCClinet,
		logger:      logger,
	}
}

func (c *WeatherGRPCClient) GetWeatherByCity(ctx context.Context, city string) (*dto.Weather, error) {
	req := &weather.GetWeatherRequest{City: city}

	resp, err := c.weatherGRPC.GetWeather(ctx, req)

	if err != nil {
		c.logger.Warnf("Failed to get weather %s:", err.Error())
		return nil, err
	}

	return mappers.MapProtoToWeatherDTO(resp), nil
}
