package httperrors

import (
	"net/http"
	domainerr "weather-forecast/internal/domain/errors"
	infraerrors "weather-forecast/internal/infrastructure/errors"
	apierrors "weather-forecast/internal/presentation/errors"
)

type HTTPError struct {
	status int
	body   map[string]any
}

func New(err error) *HTTPError {
	httpErr := &HTTPError{}
	switch err {

	case domainerr.ErrAlreadySubscribed:
		httpErr.status = http.StatusConflict

	case infraerrors.ErrInternal,
		domainerr.ErrInvalidFrequency,
		infraerrors.ErrDatabase,
		infraerrors.ErrGetWeather:
		httpErr.status = http.StatusInternalServerError

	case domainerr.ErrTokenNotFound,
		infraerrors.ErrCityNotFound:
		httpErr.status = http.StatusNotFound

	case domainerr.ErrInvalidToken,
		apierrors.ErrInvalidRequest:
		httpErr.status = http.StatusBadRequest

	default:
		httpErr.status = http.StatusInternalServerError
		httpErr.setBody("internal server error")
		return httpErr
	}
	httpErr.setBody(err.Error())
	return httpErr
}

func (err *HTTPError) Status() int {
	return err.status
}

func (err *HTTPError) Body() map[string]any {
	return err.body
}

func (err *HTTPError) setBody(msg string) {
	err.body = map[string]any{"error": msg}
}
