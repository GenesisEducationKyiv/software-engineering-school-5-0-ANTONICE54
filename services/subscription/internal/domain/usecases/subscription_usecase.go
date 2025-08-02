package usecases

import (
	"context"
	"subscription-service/internal/domain/contracts"
	domainerrors "subscription-service/internal/domain/errors"
	"subscription-service/internal/domain/models"
	"weather-forecast/pkg/logger"
)

type (
	SubscriptionRepository interface {
		Create(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
		GetByEmail(ctx context.Context, email string) (*models.Subscription, error)
		GetByToken(ctx context.Context, token string) (*models.Subscription, error)
		Update(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
		ListConfirmedByFrequency(ctx context.Context, frequency models.Frequency, lastID, pageSize int) ([]models.Subscription, error)
		DeleteByToken(ctx context.Context, token string) error
	}

	TokenManager interface {
		Generate(ctx context.Context) string
		Validate(ctx context.Context, token string) bool
	}

	NotificationSender interface {
		SendConfirmation(ctx context.Context, info *contracts.ConfirmationInfo)
		SendConfirmed(ctx context.Context, info *contracts.ConfirmedInfo)
		SendUnsubscribed(ctx context.Context, info *contracts.UnsubscribeInfo)
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

func (s *SubscriptionService) Subscribe(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {

	receivedSubsc, err := s.subscriptionRepository.GetByEmail(ctx, subscription.Email)
	if err != nil {
		return nil, err
	}
	if receivedSubsc != nil {
		return nil, domainerrors.ErrAlreadySubscribed
	}

	token := s.tokenManager.Generate(ctx)
	subscription.Token = token

	createdSubscription, err := s.subscriptionRepository.Create(ctx, *subscription)
	if err != nil {
		return nil, err
	}

	confirmationInfo := contracts.ConfirmationInfo{
		Email:     createdSubscription.Email,
		Token:     createdSubscription.Token,
		Frequency: createdSubscription.Frequency,
	}
	s.mailer.SendConfirmation(ctx, &confirmationInfo)

	return createdSubscription, nil
}

func (s *SubscriptionService) Confirm(ctx context.Context, token string) error {
	tokenIsValid := s.tokenManager.Validate(ctx, token)
	if !tokenIsValid {
		return domainerrors.ErrInvalidToken
	}

	receivedSubsc, err := s.subscriptionRepository.GetByToken(ctx, token)
	if err != nil {
		return err
	}
	if receivedSubsc == nil {
		return domainerrors.ErrTokenNotFound
	}

	if !receivedSubsc.Confirmed {
		receivedSubsc.Confirmed = true
		updatedSubsc, err := s.subscriptionRepository.Update(ctx, *receivedSubsc)
		if err != nil {
			return err
		}

		confirmedInfo := contracts.ConfirmedInfo{
			Email:     updatedSubsc.Email,
			Token:     updatedSubsc.Token,
			Frequency: updatedSubsc.Frequency,
		}

		s.mailer.SendConfirmed(ctx, &confirmedInfo)
	}

	return nil
}

func (s *SubscriptionService) Unsubscribe(ctx context.Context, token string) error {
	tokenIsValid := s.tokenManager.Validate(ctx, token)
	if !tokenIsValid {
		return domainerrors.ErrInvalidToken
	}

	receivedSubsc, err := s.subscriptionRepository.GetByToken(ctx, token)
	if err != nil {
		return err
	}
	if receivedSubsc == nil {
		return domainerrors.ErrTokenNotFound
	}

	err = s.subscriptionRepository.DeleteByToken(ctx, token)
	if err != nil {
		return err
	}

	unsubscribeInfo := contracts.UnsubscribeInfo{
		Email:     receivedSubsc.Email,
		City:      receivedSubsc.City,
		Frequency: receivedSubsc.Frequency,
	}

	s.mailer.SendUnsubscribed(ctx, &unsubscribeInfo)

	return nil
}

func (s *SubscriptionService) ListByFrequency(ctx context.Context, query *models.ListSubscriptionsQuery) ([]models.Subscription, error) {
	receivedSubscriptions, err := s.subscriptionRepository.ListConfirmedByFrequency(ctx, query.Frequency, query.LastID, query.PageSize)
	if err != nil {
		return nil, err
	}

	return receivedSubscriptions, nil
}
