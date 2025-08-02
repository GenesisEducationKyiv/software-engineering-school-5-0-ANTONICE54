package clients

import (
	"context"
	"weather-forecast/gateway/internal/mappers"
	"weather-forecast/gateway/internal/server/handlers"
	"weather-forecast/pkg/ctxutil"
	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/proto/subscription"

	"google.golang.org/grpc/metadata"
)

type (
	SubscriptionGRPCClient struct {
		subscriptionGRPC subscription.SubscriptionServiceClient
		logger           logger.Logger
	}
)

func NewSubscriptionGRPCClient(subscriptionGRPC subscription.SubscriptionServiceClient, logger logger.Logger) *SubscriptionGRPCClient {
	return &SubscriptionGRPCClient{
		subscriptionGRPC: subscriptionGRPC,
		logger:           logger,
	}
}

func (c *SubscriptionGRPCClient) Subscribe(ctx context.Context, info handlers.SubscribeRequest) error {
	processID := ctxutil.GetProcessID(ctx)
	log := c.logger.WithContext(ctx)

	md := metadata.Pairs(ctxutil.ProcessIDKey.String(), processID)
	ctx = metadata.NewOutgoingContext(ctx, md)

	log.Debugf("Calling subscribe via GRPC: Email: %s, City: %s, Frequency: %s", info.Email, info.City, info.Frequency)

	req := &subscription.SubscribeRequest{
		Email:     info.Email,
		City:      info.City,
		Frequency: mappers.MapFrequencyToProto(info.Frequency),
	}
	_, err := c.subscriptionGRPC.Subscribe(ctx, req)
	if err != nil {
		log.Warnf("Failed to subscribe via GRPC: Email: %s", info.Email)
		return err
	}

	log.Debugf("Successfully subscribed via gRPC: Email: %s", info.Email)

	return nil

}

func (c *SubscriptionGRPCClient) Confirm(ctx context.Context, token string) error {
	processID := ctxutil.GetProcessID(ctx)
	log := c.logger.WithContext(ctx)

	md := metadata.Pairs(ctxutil.ProcessIDKey.String(), processID)
	ctx = metadata.NewOutgoingContext(ctx, md)

	log.Debugf("Calling confirm via GRPC: Token: %s", token)

	req := &subscription.ConfirmRequest{
		Token: token,
	}
	_, err := c.subscriptionGRPC.Confirm(ctx, req)
	if err != nil {
		log.Warnf("Failed to confirm subscription via GRPC: Token: %s", token)
		return err
	}

	log.Debugf("Successfully confirmed subscription via gRPC: Token: %s", token)

	return nil
}

func (c *SubscriptionGRPCClient) Unsubscribe(ctx context.Context, token string) error {
	processID := ctxutil.GetProcessID(ctx)
	log := c.logger.WithContext(ctx)

	md := metadata.Pairs(ctxutil.ProcessIDKey.String(), processID)
	ctx = metadata.NewOutgoingContext(ctx, md)

	log.Debugf("Calling unsubscribe via GRPC: Token: %s", token)

	req := &subscription.UnsubscribeRequest{
		Token: token,
	}
	_, err := c.subscriptionGRPC.Unsubscribe(ctx, req)
	if err != nil {
		log.Warnf("Failed to cancel subscription via GRPC: Token: %s", token)
		return err
	}

	log.Debugf("Successfully canceled subscription via gRPC: Token: %s", token)

	return nil
}
