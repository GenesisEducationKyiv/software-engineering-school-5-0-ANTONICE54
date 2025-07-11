package errors

import "weather-forecast/pkg/apperrors"

type BroadcastServiceErrorCode string

func (c BroadcastServiceErrorCode) String() string {
	return string(c)
}

const (
	WeatherServiceErrorCode      BroadcastServiceErrorCode = "WEATHER_SERVICE_ERROR"
	SubscriptionServiceErrorCode BroadcastServiceErrorCode = "SUBSCRIPTION_SERVICE_ERROR"
	InternalErrorCode            BroadcastServiceErrorCode = "INTERNAL_ERROR"
)

var (
	InternalError = apperrors.NewInternal(InternalErrorCode, "internal server error")
)
