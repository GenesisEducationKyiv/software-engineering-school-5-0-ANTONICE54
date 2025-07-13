package mappers

import (
	"weather-broadcast-service/internal/dto"
	"weather-broadcast-service/internal/models"
	"weather-forecast/pkg/events"
	"weather-forecast/pkg/proto/subscription"
	"weather-forecast/pkg/proto/weather"
)

func MapFrequencyToProto(freq models.Frequency) subscription.Frequency {
	switch freq {
	case models.Daily:
		return subscription.Frequency_DAILY
	case models.Hourly:
		return subscription.Frequency_HOURLY
	default:
		return subscription.Frequency_DAILY
	}
}

func MapProtoToSubscriptionList(protoList *subscription.GetSubscriptionsByFrequencyResponse) *dto.SubscriptionList {
	res := &dto.SubscriptionList{
		LastIndex:     int(protoList.NextPageIndex),
		Subscriptions: make([]dto.Subscription, 0),
	}

	for _, protoSubsc := range protoList.Subscriptions {
		subsc := dto.Subscription{
			Email: protoSubsc.Email,
			City:  protoSubsc.City,
		}
		res.Subscriptions = append(res.Subscriptions, subsc)

	}
	return res

}

func MapProtoToWeatherDTO(weatherResponse *weather.GetWeatherResponse) *dto.Weather {
	return &dto.Weather{
		Temperature: weatherResponse.Temperature,
		Humidity:    int(weatherResponse.Humidity),
		Description: weatherResponse.Description,
	}
}

func MapWeatherSuccessMailToEvent(weatherInfo *dto.WeatherMailSuccessInfo) *events.WeatherSuccessEvent {
	return &events.WeatherSuccessEvent{
		Email: weatherInfo.Email,
		City:  weatherInfo.City,
		Weather: events.Weather{
			Temperature: weatherInfo.Weather.Temperature,
			Humidity:    weatherInfo.Weather.Humidity,
			Description: weatherInfo.Weather.Description,
		},
	}
}

func MapWeatherErrorMailToEvent(weatherInfo *dto.WeatherMailErrorInfo) *events.WeatherErrorEvent {
	return &events.WeatherErrorEvent{
		Email: weatherInfo.Email,
		City:  weatherInfo.City,
	}
}
