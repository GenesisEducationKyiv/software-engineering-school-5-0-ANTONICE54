package models

import (
	"weather-forecast/internal/infrastructure/apperrors"
)

type (
	Frequency string

	Subscription struct {
		ID        int
		Email     string
		City      string
		Token     string
		Frequency Frequency
		Confirmed bool
	}
)

const (
	Daily  Frequency = "daily"
	Hourly Frequency = "hourly"
)

func NewSubscription(email, city, token, frequency string) (*Subscription, error) {
	var freq Frequency
	switch frequency {
	case string(Daily):
		freq = Daily
	case string(Hourly):
		freq = Hourly
	default:
		return nil, apperrors.InvalidFrequencyInternalError
	}

	return &Subscription{
		Email:     email,
		City:      city,
		Token:     token,
		Frequency: freq,
		Confirmed: false,
	}, nil
}
