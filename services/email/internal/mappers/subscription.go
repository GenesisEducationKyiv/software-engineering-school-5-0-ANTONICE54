package mappers

import (
	"email-service/internal/dto"
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
