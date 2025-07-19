package errors

import "weather-forecast/pkg/apperrors"

type APIErrorCode string

func (c APIErrorCode) String() string {
	return string(c)
}

const (
	BadRequestErrorCode APIErrorCode = "BAD_REQUEST"
)

var (
	BadRequestError = apperrors.NewBadRequest(BadRequestErrorCode, "invalid request")
)
