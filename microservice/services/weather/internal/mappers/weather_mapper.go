package mappers

import (
	"weather-forecast/pkg/proto/weather"
	"weather-service/internal/dto"
)

func MapWeatherToProto(weatherDTO *dto.Weather) *weather.GetWeatherResponse {

	return &weather.GetWeatherResponse{
		Temperature: weatherDTO.Temperature,
		Humidity:    int32(weatherDTO.Humidity),
		Description: weatherDTO.Description,
	}
}
