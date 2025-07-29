package clients

import (
	"context"
	"weather-forecast/gateway/internal/dto"
	"weather-forecast/gateway/internal/errors"
	"weather-forecast/gateway/internal/mappers"
	"weather-forecast/pkg/apperrors"
	"weather-forecast/pkg/ctxutil"
	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/proto/weather"

	"google.golang.org/grpc/metadata"
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
	processID := ctxutil.GetProcessID(ctx)
	log := c.logger.WithField("process_id", processID)

	md := metadata.Pairs("process-id", processID)
	ctx = metadata.NewOutgoingContext(ctx, md)

	log.Debugf("Calling get weather via GRPC: %s", city)
	req := &weather.GetWeatherRequest{City: city}
	resp, err := c.weatherGRPC.GetWeather(ctx, req)
	if err != nil {
		log.Warnf("Failed to get weather via GRPC: City: %s", city)
		return nil, apperrors.FromGRPCError(err, errors.WeatherServiceErrorCode)
	}

	log.Debugf("Successfully received weather via gRPC: City %s", city)

	return mappers.MapProtoToWeatherDTO(resp), nil
}
