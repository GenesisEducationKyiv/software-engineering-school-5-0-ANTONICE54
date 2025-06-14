package handlers

import (
	"context"
	"net/http"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/apperrors"
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

	getWeatherRequest struct {
		City string `json:"city" binding:"required,alpha"`
	}
	getWeatherResponse struct {
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
	var req getWeatherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Failed to unmarshal request: %s", err.Error())
		ctx.JSON(apperrors.InvalidRequestError.Status(), apperrors.InvalidRequestError.JSON())
		return
	}

	weather, err := h.weatherService.GetWeatherByCity(ctx, req.City)

	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			ctx.JSON(appErr.Status(), appErr.JSON())
			return
		}
		ctx.JSON(apperrors.InternalError.Status(), apperrors.InternalError.JSON())
		return
	}

	response := getWeatherResponse{
		Temperature: weather.Temperature,
		Humidity:    weather.Humidity,
		Description: weather.Description,
	}

	ctx.JSON(http.StatusOK, response)

}
