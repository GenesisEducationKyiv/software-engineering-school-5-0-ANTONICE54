package handlers

import (
	"context"
	"weather-forecast/pkg/apperrors"
	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/proto/weather"
	"weather-service/internal/dto"
	"weather-service/internal/mappers"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	WeatherService interface {
		GetWeatherByCity(ctx context.Context, city string) (*dto.Weather, error)
	}

	WeatherHandler struct {
		weather.UnimplementedWeatherServiceServer
		weatherService WeatherService
		logger         logger.Logger
	}
)

func NewWeatherHandler(weatherService WeatherService, logger logger.Logger) *WeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
		logger:         logger,
	}
}

func (h *WeatherHandler) GetWeather(ctx context.Context, req *weather.GetWeatherRequest) (*weather.GetWeatherResponse, error) {

	weather, err := h.weatherService.GetWeatherByCity(ctx, req.City)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return nil, appErr.ToGRPCStatus()
		}
		return nil, status.Errorf(codes.Internal, "failed to get weather: %v", err)
	}

	return mappers.MapWeatherToProto(weather), nil
}
