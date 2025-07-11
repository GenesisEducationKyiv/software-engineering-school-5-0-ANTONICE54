package clients

import (
	"context"
	"weather-forecast/gateway/internal/errors"
	"weather-forecast/gateway/internal/mappers"
	"weather-forecast/pkg/apperrors"
	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/proto/subscription"
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

func (c *SubscriptionGRPCClient) Subscribe(ctx context.Context, email, frequency, city string) error {

	req := &subscription.SubscribeRequest{
		Email:     email,
		City:      city,
		Frequency: mappers.MapStringFrequencyToProto(frequency),
	}

	_, err := c.subscriptionGRPC.Subscribe(ctx, req)

	if err != nil {
		return apperrors.FromGRPCError(err, errors.SubscriptionServiceErrorCode)
	}
	return nil

}
func (c *SubscriptionGRPCClient) Confirm(ctx context.Context, token string) error {
	req := &subscription.ConfirmRequest{
		Token: token,
	}

	_, err := c.subscriptionGRPC.Confirm(ctx, req)

	if err != nil {
		return apperrors.FromGRPCError(err, errors.SubscriptionServiceErrorCode)
	}
	return nil
}
func (c *SubscriptionGRPCClient) Unsubscribe(ctx context.Context, token string) error {
	req := &subscription.UnsubscribeRequest{
		Token: token,
	}

	_, err := c.subscriptionGRPC.Unsubscribe(ctx, req)

	if err != nil {
		return apperrors.FromGRPCError(err, errors.SubscriptionServiceErrorCode)
	}
	return nil
}
