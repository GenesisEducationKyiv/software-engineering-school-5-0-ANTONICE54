package models

import (
	domainerr "weather-forecast/internal/domain/errors"
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
		return nil, domainerr.ErrInvalidFrequency
	}

	return &Subscription{
		Email:     email,
		City:      city,
		Token:     token,
		Frequency: freq,
		Confirmed: false,
	}, nil
}
