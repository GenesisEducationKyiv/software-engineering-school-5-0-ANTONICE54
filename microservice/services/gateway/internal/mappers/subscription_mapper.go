package mappers

import (
	"weather-forecast/gateway/internal/dto"
	"weather-forecast/pkg/proto/subscription"
)

func MapFrequencyToProto(freq dto.Frequency) subscription.Frequency {
	switch freq {
	case dto.Daily:
		return subscription.Frequency_DAILY
	case dto.Hourly:
		return subscription.Frequency_HOURLY
	default:
		return subscription.Frequency_DAILY
	}
}

func MapStringFrequencyToProto(freq string) subscription.Frequency {
	switch freq {
	case "daily":
		return subscription.Frequency_DAILY
	case "hourly":
		return subscription.Frequency_HOURLY
	default:
		return subscription.Frequency_DAILY
	}
}

func MapSubscriptionDTOToProto(subsc dto.Subscription) *subscription.SubscribeRequest {
	return &subscription.SubscribeRequest{
		Email:     subsc.Email,
		City:      subsc.City,
		Frequency: MapFrequencyToProto(subsc.Frequency),
	}
}
