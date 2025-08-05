package ctxutil

import (
	"context"
)

type contextKey string

const CorrelationIDKey contextKey = "correlation_id"

func (k contextKey) String() string {
	return string(k)
}

func GetCorrelationID(ctx context.Context) string {
	if val := ctx.Value(CorrelationIDKey.String()); val != nil {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return "unknown-process"
}
