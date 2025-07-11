package models

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
