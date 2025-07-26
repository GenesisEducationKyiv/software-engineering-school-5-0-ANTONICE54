package mappers

import (
	"net/http"
	"weather-forecast/pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HTTPFromGRPCError(err error, logger logger.Logger) (int, map[string]any) {
	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError, map[string]any{"error": "unexpected error"}
	}
	switch st.Code() {
	case codes.AlreadyExists:
		return http.StatusConflict, map[string]any{"error": st.Message()}
	case codes.InvalidArgument:
		return http.StatusBadRequest, map[string]any{"error": st.Message()}
	case codes.NotFound:
		return http.StatusNotFound, map[string]any{"error": st.Message()}
	default:
		logger.Warnf("Unexpected gRPC error: %s", err.Error())
		return http.StatusInternalServerError, map[string]any{"error": "internal server error"}
	}
}
