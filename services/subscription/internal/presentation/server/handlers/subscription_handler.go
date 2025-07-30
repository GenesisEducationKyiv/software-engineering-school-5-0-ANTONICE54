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

	subscription := mappers.SubscribeRequestToSubscribe(req)

	_, err := h.subscriptionUsecase.Subscribe(ctx, subscription)

	if err != nil {
		grpcErr := h.handleSubscribeError(err)
		return &emptypb.Empty{}, grpcErr
	}

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

	err := h.subscriptionUsecase.Confirm(ctx, req.Token)

	if err != nil {
		grpcErr := h.handleConfirmError(err)
		return &emptypb.Empty{}, grpcErr
	}

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
	err := h.subscriptionUsecase.Unsubscribe(ctx, req.Token)

	if err != nil {
		grpcErr := h.handleUnsubscribeError(err)
		return &emptypb.Empty{}, grpcErr
	}

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

	query := mappers.ProtoToListQuery(req)

	subscriptions, err := h.subscriptionUsecase.ListByFrequency(ctx, query)
	if err != nil {
		grpcErr := h.handleSubscriptionsByFrequencyError(err)
		return nil, grpcErr
	}

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
