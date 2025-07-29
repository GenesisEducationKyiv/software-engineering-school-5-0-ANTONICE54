package handlers

import (
	"context"
	"net/http"
	"weather-forecast/gateway/internal/dto"
	"weather-forecast/gateway/internal/errors"
	httperrors "weather-forecast/gateway/internal/server/http_errors"
	"weather-forecast/pkg/apperrors"
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
		httpErr := httperrors.New(errors.InvalidRequestError)
		ctx.JSON(httpErr.Status(), httpErr.JSON())
		return
	}

	log.Infof("Incoming get weather request: City: %s", req.City)

	weather, err := h.weatherClient.GetWeatherByCity(ctx, req.City)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			log.Debugf("Get weather failed for city %s: %s", req.City, appErr.Error())

			httpErr := httperrors.New(appErr)
			ctx.JSON(httpErr.Status(), httpErr.JSON())
			return
		}
		log.Errorf("Unexpected error during get weather request: %s", err.Error())

		httpErr := httperrors.New(apperrors.InternalServerError)
		ctx.JSON(httpErr.Status(), httpErr.JSON())
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
