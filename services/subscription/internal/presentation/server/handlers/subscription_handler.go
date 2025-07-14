package handlers

import (
	"context"
	"subscription-service/internal/domain/models"
	"subscription-service/internal/presentation/mappers"
	"weather-forecast/pkg/apperrors"
	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/proto/subscription"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type (
	SubscriptionUsecase interface {
		Subscribe(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error)
		Confirm(ctx context.Context, token string) error
		Unsubscribe(ctx context.Context, token string) error
		ListByFrequency(ctx context.Context, query *models.ListSubscriptionsQuery) ([]models.Subscription, error)
	}

	SubscriptionHandler struct {
		subscription.UnimplementedSubscriptionServiceServer
		subscriptionUsecase SubscriptionUsecase
		logger              logger.Logger
	}
)

func NewSubscriptionHandler(subscriptionUsecase SubscriptionUsecase, logger logger.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionUsecase: subscriptionUsecase,
		logger:              logger,
	}
}

func (h *SubscriptionHandler) Subscribe(ctx context.Context, req *subscription.SubscribeRequest) (*emptypb.Empty, error) {

	subscription := mappers.SubscribeRequestToSubscribe(req)

	_, err := h.subscriptionUsecase.Subscribe(ctx, subscription)

	if err != nil {

		if appErr, ok := err.(*apperrors.AppError); ok {
			return &emptypb.Empty{}, appErr.ToGRPCStatus()
		}
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to create subscription: %v", err)
	}

	return &emptypb.Empty{}, nil

}

func (h *SubscriptionHandler) Confirm(ctx context.Context, req *subscription.ConfirmRequest) (*emptypb.Empty, error) {

	err := h.subscriptionUsecase.Confirm(ctx, req.Token)

	if err != nil {

		if appErr, ok := err.(*apperrors.AppError); ok {
			return &emptypb.Empty{}, appErr.ToGRPCStatus()
		}
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to confirm subscription: %v", err)
	}

	return &emptypb.Empty{}, nil

}

func (h *SubscriptionHandler) Unsubscribe(ctx context.Context, req *subscription.UnsubscribeRequest) (*emptypb.Empty, error) {
	err := h.subscriptionUsecase.Unsubscribe(ctx, req.Token)

	if err != nil {

		if appErr, ok := err.(*apperrors.AppError); ok {
			return &emptypb.Empty{}, appErr.ToGRPCStatus()
		}
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to unsibscribe subscription: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *SubscriptionHandler) GetSubscriptionsByFrequency(ctx context.Context, req *subscription.GetSubscriptionsByFrequencyRequest) (*subscription.GetSubscriptionsByFrequencyResponse, error) {

	query := mappers.ProtoToListQuery(req)

	subscriptions, err := h.subscriptionUsecase.ListByFrequency(ctx, query)

	if err != nil {

		if appErr, ok := err.(*apperrors.AppError); ok {
			return nil, appErr.ToGRPCStatus()
		}
		return nil, status.Errorf(codes.Internal, "failed to list subscriptions: %v", err)
	}

	return mappers.SubscriptionListToProto(subscriptions), nil

}
