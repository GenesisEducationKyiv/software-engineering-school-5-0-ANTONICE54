package models

type (
	ListSubscriptionsQuery struct {
		Frequency Frequency
		LastID    int
		PageSize  int
	}
)
