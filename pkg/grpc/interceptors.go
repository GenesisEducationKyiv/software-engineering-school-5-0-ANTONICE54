package grpc

import (
	"context"
	"weather-forecast/pkg/ctxutil"
	"weather-forecast/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func CorrelationIDServerInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			if ids := md.Get(ctxutil.CorrelationIDKey.String()); len(ids) > 0 {
				correlationID := ids[0]
				//nolint:staticcheck
				ctx = context.WithValue(ctx, ctxutil.CorrelationIDKey.String(), correlationID)

			} else {
				log.Warnf("correlation-id not found in context")
			}
		}
		return handler(ctx, req)
	}
}

func CorrelationIDClientInterceptor(log logger.Logger) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		correlationID := ctxutil.GetCorrelationID(ctx)
		if correlationID != "" {
			md := metadata.Pairs(ctxutil.CorrelationIDKey.String(), correlationID)
			ctx = metadata.NewOutgoingContext(ctx, md)
		} else {
			log.Warnf("correlation-id not found in context")
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
