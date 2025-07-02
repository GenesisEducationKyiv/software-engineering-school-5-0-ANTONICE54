package errors

import "weather-forecast/pkg/apperrors"

type PresenantionErrorCode string

func (c PresenantionErrorCode) String() string {
	return string(c)
}

const (
	InvalidRequestErrorCode PresenantionErrorCode = "INVALID_REQUEST_ERROR"
)

var (
	InvalidRequestError = apperrors.NewBadRequest(InvalidRequestErrorCode, "invalid request")
)
