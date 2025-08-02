package handlers

import (
	"context"
	"errors"
	domainerr "subscription-service/internal/domain/errors"
	"subscription-service/internal/domain/models"
	infraerror "subscription-service/internal/infrastructure/errors"
	"subscription-service/internal/presentation/mappers"
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
		log.Warnf("Subscribe error: %s", err.Error())
		grpcErr := h.handleSubscribeError(err)
		return &emptypb.Empty{}, grpcErr
	}

	log.Infof("Subscription created successfully: id=%d, token=%s", result.ID, result.Token)
	return &emptypb.Empty{}, nil

}

func (h *SubscriptionHandler) handleSubscribeError(err error) error {
	switch {
	case errors.Is(err, domainerr.ErrAlreadySubscribed):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, infraerror.ErrDatabase):
		h.logger.Warnf("Database error during subscription: %v", err)
		return status.Error(codes.Internal, "internal server error")

	case errors.Is(err, infraerror.ErrInternal):
		return status.Error(codes.Internal, err.Error())

	default:
		h.logger.Warnf("Unexpected error during subscription: %v", err)
		return status.Error(codes.Internal, "internal server error")

	}
}

func (h *SubscriptionHandler) Confirm(ctx context.Context, req *subscription.ConfirmRequest) (*emptypb.Empty, error) {
	log := h.logger.WithContext(ctx)

	log.Infof("GRPC Confirm called: token=%s", req.Token)

	err := h.subscriptionUsecase.Confirm(ctx, req.Token)

	if err != nil {
		log.Warnf("Confirm  error for token %s: %s", req.Token, err.Error())
		grpcErr := h.handleConfirmError(err)
		return &emptypb.Empty{}, grpcErr
	}
	log.Infof("Subscription confirmed successfully: token=%s", req.Token)
	return &emptypb.Empty{}, nil

}

func (h *SubscriptionHandler) handleConfirmError(err error) error {
	switch {
	case errors.Is(err, domainerr.ErrTokenNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, domainerr.ErrInvalidToken):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, infraerror.ErrDatabase):
		h.logger.Warnf("Database error during subscription: %v", err)
		return status.Error(codes.Internal, "internal server error")

	case errors.Is(err, infraerror.ErrInternal):
		return status.Error(codes.Internal, err.Error())

	default:
		h.logger.Warnf("Unexpected error during confirmation: %v", err)
		return status.Error(codes.Internal, "internal server error")
	}

}

func (h *SubscriptionHandler) Unsubscribe(ctx context.Context, req *subscription.UnsubscribeRequest) (*emptypb.Empty, error) {
	log := h.logger.WithContext(ctx)

	log.Infof("GRPC Unsubscribe called: token=%s", req.Token)

	err := h.subscriptionUsecase.Unsubscribe(ctx, req.Token)

	if err != nil {
		log.Warnf("Unsubscribe error for token %s: %s", req.Token, err.Error())
		grpcErr := h.handleUnsubscribeError(err)
		return &emptypb.Empty{}, grpcErr
	}

	log.Infof("Subscription deleted successfully: token=%s", req.Token)

	return &emptypb.Empty{}, nil
}

func (h *SubscriptionHandler) handleUnsubscribeError(err error) error {
	switch {
	case errors.Is(err, domainerr.ErrTokenNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, domainerr.ErrInvalidToken):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, infraerror.ErrDatabase):
		h.logger.Warnf("Database error during subscription: %v", err)
		return status.Error(codes.Internal, "internal server error")

	case errors.Is(err, infraerror.ErrInternal):
		return status.Error(codes.Internal, err.Error())

	default:
		h.logger.Warnf("Unexpected error during canceling subscription: %v", err)
		return status.Error(codes.Internal, "internal server error")
	}

}

func (h *SubscriptionHandler) GetSubscriptionsByFrequency(ctx context.Context, req *subscription.GetSubscriptionsByFrequencyRequest) (*subscription.GetSubscriptionsByFrequencyResponse, error) {
	log := h.logger.WithContext(ctx)

	log.Infof("GRPC GetSubscriptionsByFrequency called: frequency=%s, lastID=%d, pageSize=%d", req.Frequency.String(), req.PageToken, req.PageSize)

	query := mappers.ProtoToListQuery(req)

	subscriptions, err := h.subscriptionUsecase.ListByFrequency(ctx, query)
	if err != nil {
		log.Warnf("ListByFrequency error: %s", err.Error())
		grpcErr := h.handleSubscriptionsByFrequencyError(err)
		return nil, grpcErr
	}

	log.Infof("Retrieved %d subscriptions for frequency %s", len(subscriptions), req.Frequency.String())

	return mappers.SubscriptionListToProto(subscriptions), nil

}

func (h *SubscriptionHandler) handleSubscriptionsByFrequencyError(err error) error {
	switch {
	case errors.Is(err, infraerror.ErrDatabase):
		h.logger.Warnf("Database error during subscription: %v", err)
		return status.Error(codes.Internal, "internal server error")

	default:
		h.logger.Warnf("Unexpected error during canceling subscription: %v", err)
		return status.Error(codes.Internal, "internal server error")
	}

}
