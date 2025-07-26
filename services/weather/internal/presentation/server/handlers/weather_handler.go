package handlers

import (
	"context"
	"errors"
	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/proto/weather"
	"weather-service/internal/domain/models"
	infraerrors "weather-service/internal/infrastructure/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	WeatherService interface {
		GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error)
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

	weatherRes, err := h.weatherService.GetWeatherByCity(ctx, req.City)
	if err != nil {
		grpcErr := h.handleGetWeatherError(err)
		return nil, grpcErr
	}

	protoWeather := &weather.GetWeatherResponse{
		Temperature: weatherRes.Temperature,
		Humidity:    int32(weatherRes.Humidity),
		Description: weatherRes.Description,
	}

	return protoWeather, nil
}

func (h *WeatherHandler) handleGetWeatherError(err error) error {

	switch {
	case errors.Is(err, infraerrors.ErrCityNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, infraerrors.ErrGetWeather):
		return status.Error(codes.Internal, err.Error())

	case errors.Is(err, infraerrors.ErrInternal):
		return status.Error(codes.Internal, err.Error())

	default:
		h.logger.Warnf("Unexpected error in GetWeather: %s", err.Error())
		return status.Error(codes.Internal, "internal server error")
	}

}
