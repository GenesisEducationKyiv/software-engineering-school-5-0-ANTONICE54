package events

type (
	Event interface {
		EventType() EventType
	}

	EventType string

	Weather struct {
		Temperature float64 `json:"temperature"`
		Humidity    int     `json:"humidity"`
		Description string  `json:"description"`
	}

	SubscriptionEvent struct {
		Email     string `json:"email"`
		Token     string `json:"token"`
		Frequency string `json:"frequency"`
	}

	ConfirmedEvent struct {
		Email     string `json:"email"`
		Token     string `json:"token"`
		Frequency string `json:"frequency"`
	}

	UnsubscribedEvent struct {
		Email     string `json:"email"`
		City      string `json:"city"`
		Frequency string `json:"frequency"`
	}

	WeatherSuccessEvent struct {
		Email   string  `json:"email"`
		City    string  `json:"city"`
		Weather Weather `json:"weather"`
	}
	WeatherErrorEvent struct {
		Email string `json:"email"`
		City  string `json:"city"`
	}
)

const (
	SubsctiptionEmail   EventType = "emails.subscription"
	ConfirmedEmail      EventType = "emails.confirmed"
	UnsubscribedEmail   EventType = "emails.unsubscribed"
	WeatherEmailSuccess EventType = "emails.weather.succes"
	WeatherEmailError   EventType = "emails.weather.error"
)

func (SubscriptionEvent) EventType() EventType {
	return SubsctiptionEmail
}

func (ConfirmedEvent) EventType() EventType {
	return ConfirmedEmail
}

func (UnsubscribedEvent) EventType() EventType {
	return UnsubscribedEmail
}

func (WeatherSuccessEvent) EventType() EventType {
	return WeatherEmailSuccess
}

func (WeatherErrorEvent) EventType() EventType {
	return WeatherEmailError
}
