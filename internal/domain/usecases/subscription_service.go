package usecases

import (
	"context"
	domainerr "weather-forecast/internal/domain/errors"
	"weather-forecast/internal/domain/models"

	"weather-forecast/internal/infrastructure/logger"
)

type (
	SubscriptionRepository interface {
		Create(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
		GetByEmail(ctx context.Context, email string) (*models.Subscription, error)
		GetByToken(ctx context.Context, token string) (*models.Subscription, error)
		Update(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
		DeleteByToken(ctx context.Context, token string) error
	}

	TokenManager interface {
		Generate(ctx context.Context) string
		Validate(ctx context.Context, token string) bool
	}

	NotificationSender interface {
		SendConfirmation(ctx context.Context, email, token string, frequency models.Frequency)
		SendConfirmed(ctx context.Context, email, token string, frequency models.Frequency)
	}

	SubscriptionService struct {
		subscriptionRepository SubscriptionRepository
		tokenManager           TokenManager
		mailer                 NotificationSender
		logger                 logger.Logger
	}
)

func NewSubscriptionService(subscriptionRepo SubscriptionRepository, tokenManager TokenManager, mailer NotificationSender, logger logger.Logger) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepository: subscriptionRepo,
		tokenManager:           tokenManager,
		mailer:                 mailer,
		logger:                 logger,
	}
}

func (s *SubscriptionService) Subscribe(ctx context.Context, email, frequency, city string) (*models.Subscription, error) {

	receivedSubsc, err := s.subscriptionRepository.GetByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	if receivedSubsc != nil {
		return nil, domainerr.ErrAlreadySubscribed
	}

	token := s.tokenManager.Generate(ctx)
	subscription, err := models.NewSubscription(email, city, token, frequency)
	if err != nil {
		return nil, err
	}

	createdSubscription, err := s.subscriptionRepository.Create(ctx, *subscription)
	if err != nil {
		return nil, err
	}

	s.mailer.SendConfirmation(ctx, createdSubscription.Email, createdSubscription.Token, createdSubscription.Frequency)

	return createdSubscription, nil
}

func (s *SubscriptionService) Confirm(ctx context.Context, token string) error {
	tokenIsValid := s.tokenManager.Validate(ctx, token)
	if !tokenIsValid {
		return domainerr.ErrInvalidToken
	}

	receivedSubsc, err := s.subscriptionRepository.GetByToken(ctx, token)
	if err != nil {
		return err
	}
	if receivedSubsc == nil {
		return domainerr.ErrTokenNotFound
	}

	if !receivedSubsc.Confirmed {
		receivedSubsc.Confirmed = true
		updatedSubsc, err := s.subscriptionRepository.Update(ctx, *receivedSubsc)

		if err != nil {
			return err
		}

		s.mailer.SendConfirmed(ctx, updatedSubsc.Email, updatedSubsc.Token, updatedSubsc.Frequency)
	}

	return nil
}

func (s *SubscriptionService) Unsubscribe(ctx context.Context, token string) error {
	tokenIsValid := s.tokenManager.Validate(ctx, token)
	if !tokenIsValid {
		return domainerr.ErrInvalidToken
	}

	receivedSubsc, err := s.subscriptionRepository.GetByToken(ctx, token)
	if err != nil {
		return err
	}
	if receivedSubsc == nil {
		return domainerr.ErrTokenNotFound
	}

	err = s.subscriptionRepository.DeleteByToken(ctx, token)
	if err != nil {
		return err
	}

	return nil
}
