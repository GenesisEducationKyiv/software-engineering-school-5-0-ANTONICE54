package errors

import (
	"net/http"
	"weather-forecast/pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	HTTPReponse struct {
		StatusCode int
		Body       map[string]any
	}
)

func NewHTTPFromGRPC(err error, logger logger.Logger) *HTTPReponse {
	st, ok := status.FromError(err)
	if !ok {
		return &HTTPReponse{
			StatusCode: http.StatusInternalServerError,
			Body:       map[string]any{"error": "unexpected error"},
		}
	}

	switch st.Code() {

	case codes.AlreadyExists:
		return &HTTPReponse{
			StatusCode: http.StatusConflict,
			Body:       map[string]any{"error": st.Message()},
		}

	case codes.InvalidArgument:
		return &HTTPReponse{
			StatusCode: http.StatusBadRequest,
			Body:       map[string]any{"error": st.Message()},
		}

	case codes.NotFound:
		return &HTTPReponse{
			StatusCode: http.StatusNotFound,
			Body:       map[string]any{"error": st.Message()},
		}

	default:
		logger.Warnf("Unexpected gRPC error: %s", err.Error())
		return &HTTPReponse{
			StatusCode: http.StatusInternalServerError,
			Body:       map[string]any{"error": "internal server error"},
		}

	}
}
