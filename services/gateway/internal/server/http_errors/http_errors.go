package httperrors

import (
	"net/http"
	"weather-forecast/pkg/apperrors"
)

type (
	HTTPError struct {
		*apperrors.AppError
	}
)

func New(err *apperrors.AppError) *HTTPError {
	return &HTTPError{
		err,
	}
}

func (err *HTTPError) Status() int {

	switch err.Type {
	case apperrors.InternalError:
		return http.StatusInternalServerError
	case apperrors.NotFoundError:
		return http.StatusNotFound
	case apperrors.ConflictError:
		return http.StatusConflict
	case apperrors.BadRequestError:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}

}

func (err *HTTPError) JSON() map[string]any {
	return map[string]any{"error": err.Message}
}
