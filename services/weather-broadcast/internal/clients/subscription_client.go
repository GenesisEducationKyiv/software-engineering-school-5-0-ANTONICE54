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
	log := c.logger.WithContext(ctx)

	log.Debugf("Calling subscription service: frequency=%s, lastID=%d, pageSize=%d", query.Frequency, query.LastID, query.PageSize)

	req := &subscription.GetSubscriptionsByFrequencyRequest{
		Frequency: mappers.MapFrequencyToProto(query.Frequency),
		PageSize:  int32(query.PageSize),
		PageToken: int32(query.LastID),
	}

	res, err := c.subscriptionGRPC.GetSubscriptionsByFrequency(ctx, req)

	if err != nil {
		log.Errorf("Subscription service call failed: %v", err)
		return nil, err
	}

	log.Infof("Retrieved %d subscriptions from service", len(res.Subscriptions))

	return mappers.MapProtoToSubscriptionList(res), nil

}
