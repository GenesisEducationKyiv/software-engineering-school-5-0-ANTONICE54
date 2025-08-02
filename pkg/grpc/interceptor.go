package grpc

import (
	"context"
	"weather-forecast/pkg/ctxutil"
	"weather-forecast/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ProcessIDInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		processID := "unknown"
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			if ids := md.Get(string(ctxutil.ProcessIDKey)); len(ids) > 0 {
				processID = ids[0]
			} else {
				log.Warnf("process-id not found in metadata")
			}
		}

		ctx = context.WithValue(ctx, ctxutil.ProcessIDKey, processID)
		return handler(ctx, req)
	}
}
