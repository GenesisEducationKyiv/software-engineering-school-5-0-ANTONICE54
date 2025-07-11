package dto

type (
	Subscription struct {
		Email     string
		City      string
		Frequency Frequency
	}

	Frequency string
)

const (
	Daily  Frequency = "daily"
	Hourly Frequency = "hourly"
)
