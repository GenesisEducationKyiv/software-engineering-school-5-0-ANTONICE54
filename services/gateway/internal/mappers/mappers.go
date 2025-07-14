package mappers

import (
	"weather-forecast/gateway/internal/dto"
	"weather-forecast/pkg/proto/subscription"
	"weather-forecast/pkg/proto/weather"
)

func MapFrequencyToProto(freq string) subscription.Frequency {
	switch freq {
	case "daily":
		return subscription.Frequency_DAILY
	case "hourly":
		return subscription.Frequency_HOURLY
	default:
		return subscription.Frequency_DAILY
	}
}

func MapProtoToWeatherDTO(weatherResponse *weather.GetWeatherResponse) *dto.Weather {
	return &dto.Weather{
		Temperature: weatherResponse.Temperature,
		Humidity:    int(weatherResponse.Humidity),
		Description: weatherResponse.Description,
	}
}
