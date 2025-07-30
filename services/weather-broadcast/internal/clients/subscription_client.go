package clients

import (
	"context"
	"weather-broadcast-service/internal/dto"
	"weather-broadcast-service/internal/mappers"

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

func (c *SubscriptionGRPCClient) ListByFrequency(ctx context.Context, query dto.ListSubscriptionsQuery) (*dto.SubscriptionList, error) {

	req := &subscription.GetSubscriptionsByFrequencyRequest{
		Frequency: mappers.MapFrequencyToProto(query.Frequency),
		PageSize:  int32(query.PageSize),
		PageToken: int32(query.LastID),
	}

	res, err := c.subscriptionGRPC.GetSubscriptionsByFrequency(ctx, req)

	if err != nil {
		c.logger.Warnf("Failed to get subscription list %s:", err.Error())
		return nil, err
	}
	return mappers.MapProtoToSubscriptionList(res), nil

}
