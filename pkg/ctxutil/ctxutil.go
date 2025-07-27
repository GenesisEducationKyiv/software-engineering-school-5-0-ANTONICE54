package ctxutil

import "context"

func GetProcessID(ctx context.Context) string {
	if val := ctx.Value("process_id"); val != nil {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return "unknown-process"
}
