package mappers

import (
	"email-service/internal/dto"
	"weather-forecast/pkg/proto/events"
)

func SubscribeEventToDTO(event *events.SubscriptionEvent) *dto.SubscriptionEmailInfo {
	return &dto.SubscriptionEmailInfo{
		Email:     event.Email,
		Frequency: event.Frequency,
		Token:     event.Token,
	}
}

func ConfirmEventToDTO(event *events.ConfirmedEvent) *dto.ConfirmedEmailInfo {
	return &dto.ConfirmedEmailInfo{
		Email:     event.Email,
		Frequency: event.Frequency,
		Token:     event.Token,
	}
}

func WeatherToDTO(weather *events.Weather) *dto.Weather {
	return &dto.Weather{
		Temperature: weather.Temperature,
		Humidity:    int(weather.Humidity),
		Description: weather.Description,
	}
}

func SuccessWeatehrToDTO(event *events.WeatherSuccessEvent) *dto.WeatherSuccess {
	weather := WeatherToDTO(event.Weather)

	return &dto.WeatherSuccess{
		City:    event.City,
		Email:   event.Email,
		Weather: *weather,
	}
}

func ErrorWeatehrToDTO(event *events.WeatherErrorEvent) *dto.WeatherError {
	return &dto.WeatherError{
		City:  event.City,
		Email: event.Email,
	}
}
