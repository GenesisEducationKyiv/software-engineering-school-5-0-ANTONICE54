package services

import (
	"context"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/apperrors"
	"weather-forecast/internal/infrastructure/logger"
)

type (
	SubscriptionUseCase interface {
		Subscribe(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
		Confirm(ctx context.Context, token string) (*models.Subscription, error)
		Unsubscribe(ctx context.Context, token string) error
	}
	TokenManager interface {
		Generate(ctx context.Context) string
		Validate(ctx context.Context, token string) bool
	}

	NotificationServiceI interface {
		SendConfirmation(ctx context.Context, email, token string, frequency models.Frequency)
		SendConfirmed(ctx context.Context, email, token string, frequency models.Frequency)
	}

	SubscriptionService struct {
		subscriptionUC SubscriptionUseCase
		tokenManager   TokenManager
		mailer         NotificationServiceI
		logger         logger.Logger
	}
)

func NewSubscriptionService(subscriptionUC SubscriptionUseCase, tokenManager TokenManager, mailer NotificationServiceI, logger logger.Logger) *SubscriptionService {
	return &SubscriptionService{
		subscriptionUC: subscriptionUC,
		tokenManager:   tokenManager,
		mailer:         mailer,
		logger:         logger,
	}
}

func (s *SubscriptionService) Subscribe(ctx context.Context, email, frequency, city string) (*models.Subscription, error) {
	token := s.tokenManager.Generate(ctx)

	subscription := models.Subscription{
		Email:     email,
		Frequency: models.Frequency(frequency),
		City:      city,
		Token:     token,
	}
	result, err := s.subscriptionUC.Subscribe(ctx, subscription)

	if err != nil {
		return nil, err
	}

	s.mailer.SendConfirmation(ctx, result.Email, result.Token, result.Frequency)

	return result, nil
}

func (s *SubscriptionService) Confirm(ctx context.Context, token string) error {
	tokenIsValid := s.tokenManager.Validate(ctx, token)

	if !tokenIsValid {
		return apperrors.InvalidTokenError
	}

	subsc, err := s.subscriptionUC.Confirm(ctx, token)
	if err != nil {
		return err
	}

	s.mailer.SendConfirmed(ctx, subsc.Email, subsc.Token, subsc.Frequency)

	return nil
}

func (s *SubscriptionService) Unsubscribe(ctx context.Context, token string) error {
	tokenIsValid := s.tokenManager.Validate(ctx, token)

	if !tokenIsValid {
		return apperrors.InvalidTokenError
	}

	err := s.subscriptionUC.Unsubscribe(ctx, token)
	if err != nil {
		return err
	}

	return nil
}
