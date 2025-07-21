package apperrors

import "net/http"

type (
	AppError struct {
		status      int
		message     string
		jsonMessage map[string]any
	}
)

func New(status int, msg string, jsonMsg string) *AppError {
	return &AppError{
		status:      status,
		message:     msg,
		jsonMessage: map[string]any{"error": jsonMsg},
	}
}

func (err *AppError) Error() string {
	return err.message
}

func (err *AppError) Status() int {
	return err.status
}

func (err *AppError) JSON() map[string]any {
	return err.jsonMessage
}

var (
	DatabaseError                 = New(http.StatusInternalServerError, "Database raised an error", "database raised an error")
	AlreadySubscribedError        = New(http.StatusConflict, "Already subscribed", "email already subscribed")
	TokenNotFoundError            = New(http.StatusNotFound, "Token not found", "there is no subscription with such token")
	CityNotFoundError             = New(http.StatusNotFound, "City not found", "there is no city with such name")
	GetWeatherError               = New(http.StatusInternalServerError, "Failed to get weather", "failed to get weather")
	InvalidRequestError           = New(http.StatusBadRequest, "Invalid request", "invalid request")
	InvalidTokenError             = New(http.StatusBadRequest, "Invalid token", "invalid token")
	InternalError                 = New(http.StatusInternalServerError, "Internal server error", "internal server error")
	InvalidFrequencyInternalError = New(http.StatusInternalServerError, "Unexpected frequency value", "internal: unexpected frequency value")
)
