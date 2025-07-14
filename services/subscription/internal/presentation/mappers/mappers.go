package mappers

import (
	"subscription-service/internal/domain/models"
	"weather-forecast/pkg/proto/subscription"
)

func ProtoToFreaquency(protoFrequency subscription.Frequency) models.Frequency {

	switch protoFrequency {
	case subscription.Frequency_DAILY:
		return models.Daily
	case subscription.Frequency_HOURLY:
		return models.Hourly
	default:
		return models.Daily
	}

}

func FrequencyToProto(freq models.Frequency) subscription.Frequency {
	switch freq {
	case models.Daily:
		return subscription.Frequency_DAILY
	case models.Hourly:
		return subscription.Frequency_HOURLY
	default:
		return subscription.Frequency_DAILY
	}
}

func SubscriptionToProto(subsc models.Subscription) *subscription.Subscription {
	return &subscription.Subscription{
		Email: subsc.Email,
		City:  subsc.City,
	}
}

func SubscriptionListToProto(subscriptions []models.Subscription) *subscription.GetSubscriptionsByFrequencyResponse {
	protoSubscList := make([]*subscription.Subscription, 0)
	lastIndex := 0

	for _, subsc := range subscriptions {
		lastIndex = subsc.ID

		protoSubsc := SubscriptionToProto(subsc)

		protoSubscList = append(protoSubscList, protoSubsc)

	}

	return &subscription.GetSubscriptionsByFrequencyResponse{
		Subscriptions: protoSubscList,
		NextPageIndex: int32(lastIndex),
	}

}

func SubscribeRequestToSubscribe(req *subscription.SubscribeRequest) *models.Subscription {
	return &models.Subscription{
		Email:     req.Email,
		City:      req.City,
		Frequency: ProtoToFreaquency(req.Frequency),
		Confirmed: false,
	}
}

func ProtoToListQuery(req *subscription.GetSubscriptionsByFrequencyRequest) *models.ListSubscriptionsQuery {
	return &models.ListSubscriptionsQuery{
		Frequency: ProtoToFreaquency(req.Frequency),
		LastID:    int(req.PageToken),
		PageSize:  int(req.PageSize),
	}
}
