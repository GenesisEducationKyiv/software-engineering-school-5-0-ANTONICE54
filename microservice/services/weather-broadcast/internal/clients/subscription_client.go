package clients

import (
	"context"
	"weather-broadcast-service/internal/dto"
	"weather-broadcast-service/internal/errors"
	"weather-broadcast-service/internal/mappers"

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

func (c *SubscriptionGRPCClient) ListByFrequency(ctx context.Context, frequency dto.Frequency, pageToken, pageSize int) (*dto.SubscriptionList, error) {

	req := &subscription.GetSubscriptionsByFrequencyRequest{
		Frequency: mappers.MapDTOFrequencyToProto(frequency),
		PageSize:  int32(pageSize),
		PageToken: int32(pageToken),
	}

	res, err := c.subscriptionGRPC.GetSubscriptionsByFrequency(ctx, req)

	if err != nil {
		return nil, apperrors.FromGRPCError(err, errors.SubscriptionServiceErrorCode)
	}
	return mappers.MapProtoToSubscriptionList(res), nil

}
