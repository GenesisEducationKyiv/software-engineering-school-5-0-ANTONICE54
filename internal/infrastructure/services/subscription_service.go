package services

import (
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/apperrors"
	"weather-forecast/internal/infrastructure/logger"
)

type (
	SubscriptionUseCaseI interface {
		Subscribe(subscription models.Subscription) (*models.Subscription, error)
		Confirm(token string) (*models.Subscription, error)
		Unsubscribe(token string) error
	}
	TokenManagerI interface {
		Generate() string
		Validate(token string) bool
	}

	NotificationServiceI interface {
		SendConfirmation(email, token string, frequency models.Frequency)
		SendConfirmed(email, token string, frequency models.Frequency)
	}

	SubscriptionService struct {
		subscriptionUC SubscriptionUseCaseI
		tokenManager   TokenManagerI
		mailer         NotificationServiceI
		logger         logger.Logger
	}
)

func NewSubscriptionService(subscriptionUC SubscriptionUseCaseI, tokenManager TokenManagerI, mailer NotificationServiceI, logger logger.Logger) *SubscriptionService {
	return &SubscriptionService{
		subscriptionUC: subscriptionUC,
		tokenManager:   tokenManager,
		mailer:         mailer,
		logger:         logger,
	}
}

func (s *SubscriptionService) Subscribe(subscription models.Subscription) (*models.Subscription, error) {
	token := s.tokenManager.Generate()
	subscription.Token = token
	result, err := s.subscriptionUC.Subscribe(subscription)

	if err != nil {
		return nil, err
	}

	s.mailer.SendConfirmation(result.Email, result.Token, result.Frequency)

	return result, nil
}

func (s *SubscriptionService) Confirm(token string) error {
	tokenIsValid := s.tokenManager.Validate(token)

	if !tokenIsValid {
		return apperrors.InvalidTokenError
	}

	subsc, err := s.subscriptionUC.Confirm(token)
	if err != nil {
		return err
	}

	s.mailer.SendConfirmed(subsc.Email, subsc.Token, subsc.Frequency)

	return nil
}

func (s *SubscriptionService) Unsubscribe(token string) error {
	tokenIsValid := s.tokenManager.Validate(token)

	if !tokenIsValid {
		return apperrors.InvalidTokenError
	}

	err := s.subscriptionUC.Unsubscribe(token)
	if err != nil {
		return err
	}

	return nil
}
