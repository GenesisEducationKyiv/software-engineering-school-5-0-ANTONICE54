package grpc

import (
	"context"
	"weather-forecast/pkg/ctxutil"
	"weather-forecast/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func CorrelationIDInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		correlationID := "unknown"
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			if ids := md.Get(ctxutil.CorrelationIDKey.String()); len(ids) > 0 {
				correlationID = ids[0]
			} else {
				log.Warnf("process-id not found in metadata")
			}
		}
		//nolint:staticcheck
		ctx = context.WithValue(ctx, ctxutil.CorrelationIDKey.String(), correlationID)
		return handler(ctx, req)
	}
}
