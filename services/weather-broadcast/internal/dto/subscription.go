package dto

import "weather-broadcast-service/internal/models"

type (
	Subscription struct {
		Email string
		City  string
	}
	SubscriptionList struct {
		Subscriptions []Subscription
		LastIndex     int
	}

	ListSubscriptionsQuery struct {
		Frequency models.Frequency
		LastID    int
		PageSize  int
	}
)
