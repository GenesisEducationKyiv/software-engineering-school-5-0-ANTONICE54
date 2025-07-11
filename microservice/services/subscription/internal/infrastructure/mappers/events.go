package mappers

import (
	"subscription-service/internal/domain/dto"
	"weather-forecast/pkg/events"
)

func ConfirmationInfoToEvent(info *dto.ConfirmationInfo) *events.SubscriptionEvent {
	return &events.SubscriptionEvent{
		Email:     info.Email,
		Token:     info.Token,
		Frequency: string(info.Frequency),
	}
}

func ConfirmedInfoToEvent(info *dto.ConfirmedInfo) *events.ConfirmedEvent {
	return &events.ConfirmedEvent{
		Email:     info.Email,
		Token:     info.Token,
		Frequency: string(info.Frequency),
	}
}
