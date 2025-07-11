package mappers

import (
	"email-service/internal/dto"
	"weather-forecast/pkg/events"
)

func WeatherToDTO(weather events.Weather) *dto.Weather {
	return &dto.Weather{
		Temperature: weather.Temperature,
		Humidity:    weather.Humidity,
		Description: weather.Description,
	}
}

func SuccessWeatehrToDTO(event events.WeatherSuccessEvent) *dto.WeatherSuccess {
	weather := WeatherToDTO(event.Weather)

	return &dto.WeatherSuccess{
		City:    event.City,
		Email:   event.Email,
		Weather: *weather,
	}
}

func ErrorWeatehrToDTO(event events.WeatherErrorEvent) *dto.WeatherError {
	return &dto.WeatherError{
		City:  event.City,
		Email: event.Email,
	}
}
