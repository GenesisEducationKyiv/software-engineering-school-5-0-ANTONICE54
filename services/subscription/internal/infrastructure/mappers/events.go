package mappers

import (
	"subscription-service/internal/domain/contracts"
	"weather-forecast/pkg/events"
)

func ConfirmationInfoToEvent(info *contracts.ConfirmationInfo) *events.SubscriptionEvent {
	return &events.SubscriptionEvent{
		Email:     info.Email,
		Token:     info.Token,
		Frequency: string(info.Frequency),
	}
}

func ConfirmedInfoToEvent(info *contracts.ConfirmedInfo) *events.ConfirmedEvent {
	return &events.ConfirmedEvent{
		Email:     info.Email,
		Token:     info.Token,
		Frequency: string(info.Frequency),
	}
}
