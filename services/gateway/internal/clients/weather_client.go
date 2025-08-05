package clients

import (
	"context"
	"weather-forecast/gateway/internal/dto"
	"weather-forecast/gateway/internal/mappers"
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

func NewWeatherGRPCClient(weatherGRPCClient weather.WeatherServiceClient, logger logger.Logger) *WeatherGRPCClient {
	return &WeatherGRPCClient{
		weatherGRPC: weatherGRPCClient,
		logger:      logger,
	}
}

func (c *WeatherGRPCClient) GetWeatherByCity(ctx context.Context, city string) (*dto.Weather, error) {
	correlationID := ctxutil.GetCorrelationID(ctx)
	log := c.logger.WithContext(ctx)

	md := metadata.Pairs(ctxutil.CorrelationIDKey.String(), correlationID)
	ctx = metadata.NewOutgoingContext(ctx, md)

	log.Debugf("Calling get weather via GRPC: %s", city)
	req := &weather.GetWeatherRequest{City: city}
	resp, err := c.weatherGRPC.GetWeather(ctx, req)
	if err != nil {
		log.Warnf("Failed to get weather via GRPC: City: %s", city)
		return nil, err
	}

	log.Debugf("Successfully received weather via gRPC: City %s", city)

	return mappers.MapProtoToWeatherDTO(resp), nil
}
