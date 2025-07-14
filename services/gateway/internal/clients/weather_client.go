package clients

import (
	"context"
	"weather-forecast/gateway/internal/dto"
	"weather-forecast/gateway/internal/errors"
	"weather-forecast/gateway/internal/mappers"
	"weather-forecast/pkg/apperrors"
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
		return nil, apperrors.FromGRPCError(err, errors.WeatherServiceErrorCode)
	}

	return mappers.MapProtoToWeatherDTO(resp), nil
}
