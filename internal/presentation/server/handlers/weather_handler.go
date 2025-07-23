package handlers

import (
	"context"
	"net/http"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/logger"
	apierrors "weather-forecast/internal/presentation/errors"
	"weather-forecast/internal/presentation/httperrors"

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
		httpErr := httperrors.New(apierrors.ErrInvalidRequest)
		ctx.JSON(httpErr.Status(), httpErr.Body())
		return
	}

	weather, err := h.weatherService.GetWeatherByCity(ctx, req.City)

	if err != nil {
		httpErr := httperrors.New(err)
		ctx.JSON(httpErr.Status(), httpErr.Body())
		return
	}

	response := GetWeatherResponse{
		Temperature: weather.Temperature,
		Humidity:    weather.Humidity,
		Description: weather.Description,
	}

	ctx.JSON(http.StatusOK, response)

}
