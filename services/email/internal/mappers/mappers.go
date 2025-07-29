package mappers

import (
	"email-service/internal/dto"
	"strings"
	"weather-forecast/pkg/events"
)

func SubscribeEventToDTO(event events.SubscriptionEvent) *dto.SubscriptionEmailInfo {
	return &dto.SubscriptionEmailInfo{
		Email:     event.Email,
		Frequency: event.Frequency,
		Token:     event.Token,
	}
}

func ConfirmEventToDTO(event events.ConfirmedEvent) *dto.ConfirmedEmailInfo {
	return &dto.ConfirmedEmailInfo{
		Email:     event.Email,
		Frequency: event.Frequency,
		Token:     event.Token,
	}
}

func UnsubscribeEventToDTO(event events.UnsubscribedEvent) *dto.UnsubscribedEmailInfo {
	return &dto.UnsubscribedEmailInfo{
		Email:     event.Email,
		City:      event.City,
		Frequency: event.Frequency,
	}
}

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

func SubjectToSubjectType(subject string) string {
	subject = strings.ToLower(subject)
	switch {
	case strings.Contains(subject, "confirm"):
		return "subscription"
	case strings.Contains(subject, "weather"):
		return "weather"
	case strings.Contains(subject, "canceled"):
		return "unsubscribe"
	case strings.Contains(subject, "comfirmed"):
		return "comfirmation"
	default:
		return "other"
	}
}
