package handlers

import (
	"context"
	"errors"
	"net/http"
	"weather-forecast/internal/domain/models"
	infraerrors "weather-forecast/internal/infrastructure/errors"
	"weather-forecast/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
)

type (
	WeatherService interface {
		GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error)
	}

	WeatherHandler struct {
		weatherService WeatherService
		logger         logger.Logger
	}

	GetWeatherRequest struct {
		City string `json:"city" binding:"required,alpha"`
	}
	GetWeatherResponse struct {
		Temperature float64 `json:"temperature"`
		Humidity    int     `json:"humidity"`
		Description string  `json:"description"`
	}
)

func NewWeatherHandler(weatherService WeatherService, logger logger.Logger) *WeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
		logger:         logger,
	}
}

func (h *WeatherHandler) Get(ctx *gin.Context) {
	var req GetWeatherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Failed to unmarshal request: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	weather, err := h.weatherService.GetWeatherByCity(ctx, req.City)

	if err != nil {
		h.handleGetError(ctx, err)
		return
	}

	response := GetWeatherResponse{
		Temperature: weather.Temperature,
		Humidity:    weather.Humidity,
		Description: weather.Description,
	}

	ctx.JSON(http.StatusOK, response)

}

func (h *WeatherHandler) handleGetError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, infraerrors.ErrCityNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

	case errors.Is(err, infraerrors.ErrGetWeather):
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	case errors.Is(err, infraerrors.ErrInternal):
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	default:
		h.logger.Warnf("Unexpected error during subscription: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

}
