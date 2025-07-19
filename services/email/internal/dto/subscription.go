package dto

type (
	ConfirmedEmailInfo struct {
		Email     string
		Token     string
		Frequency string
	}
	SubscriptionEmailInfo struct {
		Email     string
		Token     string
		Frequency string
	}
	UnsubscribedEmailInfo struct {
		Email     string
		City      string
		Frequency string
	}
)
