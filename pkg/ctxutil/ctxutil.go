package ctxutil

import (
	"context"
)

type contextKey string

const ProcessIDKey contextKey = "process_id"

func (k contextKey) String() string {
	return string(k)
}

func GetProcessID(ctx context.Context) string {
	if val := ctx.Value(ProcessIDKey.String()); val != nil {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return "unknown-process"
}
