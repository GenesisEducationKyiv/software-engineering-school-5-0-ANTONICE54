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
	log := h.logger.WithContext(ctx)

	log.Infof("GRPC Subscribe called: email=%s, city=%s, frequency=%s", req.Email, req.City, req.Frequency.String())

	subscription := mappers.SubscribeRequestToSubscribe(req)

	result, err := h.subscriptionUsecase.Subscribe(ctx, subscription)

	if err != nil {

		if appErr, ok := err.(*apperrors.AppError); ok {
			log.Warnf("Subscribe error: %s", appErr.Message)
			return &emptypb.Empty{}, appErr.ToGRPCStatus()
		}
		log.Errorf("Subscribe unexpected error: %v", err)
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to create subscription: %v", err)
	}

	log.Infof("Subscription created successfully: id=%d, token=%s", result.ID, result.Token)
	return &emptypb.Empty{}, nil

}

func (h *SubscriptionHandler) Confirm(ctx context.Context, req *subscription.ConfirmRequest) (*emptypb.Empty, error) {
	log := h.logger.WithContext(ctx)

	log.Infof("GRPC Confirm called: token=%s", req.Token)

	err := h.subscriptionUsecase.Confirm(ctx, req.Token)

	if err != nil {

		if appErr, ok := err.(*apperrors.AppError); ok {
			log.Warnf("Confirm  error for token %s: %s", req.Token, appErr.Message)
			return &emptypb.Empty{}, appErr.ToGRPCStatus()
		}
		log.Errorf("Confirm unexpected error for token %s: %v", req.Token, err)
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to confirm subscription: %v", err)
	}
	log.Infof("Subscription confirmed successfully: token=%s", req.Token)
	return &emptypb.Empty{}, nil

}

func (h *SubscriptionHandler) Unsubscribe(ctx context.Context, req *subscription.UnsubscribeRequest) (*emptypb.Empty, error) {
	log := h.logger.WithContext(ctx)

	log.Infof("GRPC Unsubscribe called: token=%s", req.Token)

	err := h.subscriptionUsecase.Unsubscribe(ctx, req.Token)

	if err != nil {

		if appErr, ok := err.(*apperrors.AppError); ok {
			log.Warnf("Unsubscribe error for token %s: %s", req.Token, appErr.Message)
			return &emptypb.Empty{}, appErr.ToGRPCStatus()
		}
		log.Errorf("Unsubscribe unexpected error for token %s: %v", req.Token, err)
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to unsibscribe subscription: %v", err)
	}

	log.Infof("Subscription deleted successfully: token=%s", req.Token)

	return &emptypb.Empty{}, nil
}

func (h *SubscriptionHandler) GetSubscriptionsByFrequency(ctx context.Context, req *subscription.GetSubscriptionsByFrequencyRequest) (*subscription.GetSubscriptionsByFrequencyResponse, error) {
	log := h.logger.WithContext(ctx)

	log.Infof("GRPC GetSubscriptionsByFrequency called: frequency=%s, lastID=%d, pageSize=%d", req.Frequency.String(), req.PageToken, req.PageSize)

	query := mappers.ProtoToListQuery(req)

	subscriptions, err := h.subscriptionUsecase.ListByFrequency(ctx, query)

	if err != nil {

		if appErr, ok := err.(*apperrors.AppError); ok {
			log.Warnf("ListByFrequency error: %s", appErr.Message)
			return nil, appErr.ToGRPCStatus()
		}
		log.Errorf("ListByFrequency unexpected error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list subscriptions: %v", err)
	}

	log.Infof("Retrieved %d subscriptions for frequency %s", len(subscriptions), req.Frequency.String())

	return mappers.SubscriptionListToProto(subscriptions), nil

}
