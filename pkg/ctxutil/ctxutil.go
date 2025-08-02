package ctxutil

import "context"

type contextKey string

const ProcessIDKey contextKey = "process_id"

func GetProcessID(ctx context.Context) string {
	if val := ctx.Value("process_id"); val != nil {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return "unknown-process"
}
