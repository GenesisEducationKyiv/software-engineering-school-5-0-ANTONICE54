package errors

type BroadcastServiceErrorCode string

func (c BroadcastServiceErrorCode) String() string {
	return string(c)
}

const (
	WeatherServiceErrorCode      BroadcastServiceErrorCode = "WEATHER_SERVICE_ERROR"
	SubscriptionServiceErrorCode BroadcastServiceErrorCode = "SUBSCRIPTION_SERVICE_ERROR"
)

var ()
