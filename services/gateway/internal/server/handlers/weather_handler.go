package handlers

import (
	"context"
	"net/http"
	"weather-forecast/gateway/internal/dto"
	"weather-forecast/gateway/internal/errors"
	"weather-forecast/pkg/logger"

	"github.com/gin-gonic/gin"
)

type (
	WeatherClient interface {
		GetWeatherByCity(ctx context.Context, city string) (*dto.Weather, error)
	}

	WeatherHandler struct {
		weatherClient WeatherClient
		logger        logger.Logger
	}

	GetWeatherRequest struct {
		City string `json:"city" binding:"required"`
	}
	GetWeatherResponse struct {
		Temperature float64 `json:"temperature"`
		Humidity    int     `json:"humidity"`
		Description string  `json:"description"`
	}
)

func NewWeatherHandler(weatherClient WeatherClient, logger logger.Logger) *WeatherHandler {
	return &WeatherHandler{
		weatherClient: weatherClient,
		logger:        logger,
	}
}

func (h *WeatherHandler) Get(ctx *gin.Context) {
	log := h.logger.WithContext(ctx)

	var req GetWeatherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Debugf("Failed to unmarshal request: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	log.Infof("Incoming get weather request: City: %s", req.City)

	weather, err := h.weatherClient.GetWeatherByCity(ctx, req.City)
	if err != nil {
		log.Debugf("Get weather failed for city %s: %s", req.City, err.Error())
		httpErr := errors.NewHTTPFromGRPC(err, h.logger)
		ctx.JSON(httpErr.StatusCode, httpErr.Body)
		return
	}

	log.Infof("Weather successfully retrieved: City: %s", req.City)

	response := GetWeatherResponse{
		Temperature: weather.Temperature,
		Humidity:    weather.Humidity,
		Description: weather.Description,
	}

	ctx.JSON(http.StatusOK, response)

}
