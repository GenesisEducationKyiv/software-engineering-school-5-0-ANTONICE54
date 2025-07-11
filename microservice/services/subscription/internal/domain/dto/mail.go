package dto

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
)
