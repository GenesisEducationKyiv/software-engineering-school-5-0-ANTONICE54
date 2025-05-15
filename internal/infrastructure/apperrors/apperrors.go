package apperrors

import "net/http"

type (
	Type string

	AppError struct {
		Type        Type
		Message     string
		JSONMessage map[string]any
	}
)

const (
	Internal   Type = "INTERNAL"
	BadRequest Type = "BAD_REQUEST"
	Conflict   Type = "CONFLICT"
	NotFound   Type = "NOT_FOUND"
)

func New(errType Type, msg string, jsonMsg string) *AppError {
	return &AppError{
		Type:        errType,
		Message:     msg,
		JSONMessage: map[string]any{"error": jsonMsg},
	}
}

func (err *AppError) Error() string {
	return err.Message
}

func (err *AppError) Status() int {
	switch err.Type {
	case BadRequest:
		return http.StatusBadRequest
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case Conflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

var (
	DatabaseError          = New(Internal, "Database raised an error", "database raised an error")
	AlreadySubscribedError = New(Conflict, "Already subscribed", "email already subscribed")
	TokenNotFoundError     = New(NotFound, "Token not found", "there is no subscription with such token")
	CityNotFoundError      = New(NotFound, "City not found", "there is no city with such name")
	GetWeatherError        = New(Internal, "Failed to get weather", "failed to get weather")
	InvalidRequestError    = New(BadRequest, "Invalid request", "invalid request")
	InvalidTokenError      = New(BadRequest, "Invalid token", "invalid token")
	InternalError          = New(Internal, "Internal server error", "internal server error")
)
