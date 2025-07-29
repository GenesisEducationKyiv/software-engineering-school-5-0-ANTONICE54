package handlers

import (
	"context"
	"weather-forecast/pkg/apperrors"
	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/proto/weather"
	"weather-service/internal/domain/models"

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
	log := h.logger.WithContext(ctx)

	log.Infof("GRPC GetWeather called: city=%s", req.City)
	weatherRes, err := h.weatherService.GetWeatherByCity(ctx, req.City)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			log.Warnf("GetWeather error: %s", appErr.Message)

			return nil, appErr.ToGRPCStatus()
		}
		log.Errorf("GetWeather unexpected error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get weather: %v", err)
	}

	protoWeather := &weather.GetWeatherResponse{
		Temperature: weatherRes.Temperature,
		Humidity:    int32(weatherRes.Humidity),
		Description: weatherRes.Description,
	}

	log.Infof("Weather recieved successfully: city=%s", req.City)

	return protoWeather, nil
}
