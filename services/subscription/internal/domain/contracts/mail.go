package contracts

import "subscription-service/internal/domain/models"

type (
	ConfirmationInfo struct {
		Email     string
		Token     string
		Frequency models.Frequency
	}

	ConfirmedInfo struct {
		Email     string
		Token     string
		Frequency models.Frequency
	}
	UnsubscribeInfo struct {
		Email     string
		City      string
		Frequency models.Frequency
	}
)
