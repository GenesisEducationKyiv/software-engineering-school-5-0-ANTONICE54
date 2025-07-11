package mappers

import (
	"weather-forecast/gateway/internal/dto"
	"weather-forecast/pkg/proto/weather"
)

func MapProtoToWeatherDTO(weatherResponse *weather.GetWeatherResponse) *dto.WeatherDTO {
	return &dto.WeatherDTO{
		Temperature: weatherResponse.Temperature,
		Humidity:    int(weatherResponse.Humidity),
		Description: weatherResponse.Description,
	}
}
