package dto

type (
	Frequency string

	Subscription struct {
		Email string
		City  string
	}
	SubscriptionList struct {
		Subscriptions []Subscription
		LastIndex     int
	}
)

const (
	Daily  Frequency = "daily"
	Hourly Frequency = "hourly"
)
