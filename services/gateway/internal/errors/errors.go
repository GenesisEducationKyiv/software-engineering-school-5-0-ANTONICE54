package errors

import "weather-forecast/pkg/apperrors"

type GatewayErrorCode string

func (c GatewayErrorCode) String() string {
	return string(c)
}

const (
	WeatherServiceErrorCode      GatewayErrorCode = "WEATHER_SERVICE_ERROR"
	SubscriptionServiceErrorCode GatewayErrorCode = "SUBSCRIPTION_SERVICE_ERROR"
	InvalidRequestErrorCode      GatewayErrorCode = "INVALID_REQUEST_ERROR"
)

var (
	InvalidRequestError = apperrors.NewBadRequest(InvalidRequestErrorCode, "bad request")
)
