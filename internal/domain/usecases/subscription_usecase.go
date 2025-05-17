package usecases

import (
	"context"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/apperrors"
	"weather-forecast/internal/infrastructure/logger"
)

type (
	SubscriptionRepository interface {
		Create(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
		GetByEmail(ctx context.Context, email string) (*models.Subscription, error)
		GetByToken(ctx context.Context, token string) (*models.Subscription, error)
		Update(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
		DeleteByToken(ctx context.Context, token string) error
		ListConfirmedByFrequency(ctx context.Context, frequency models.Frequency, lastID, pageSize int) ([]models.Subscription, error)
	}
	SubscriptionUseCase struct {
		subscriptionRepository SubscriptionRepository
		logger                 logger.Logger
	}
)

func NewSubscriptionUseCase(subscRepo SubscriptionRepository, logger logger.Logger) *SubscriptionUseCase {
	return &SubscriptionUseCase{
		subscriptionRepository: subscRepo,
		logger:                 logger,
	}
}

func (u *SubscriptionUseCase) Subscribe(ctx context.Context, subscription models.Subscription) (*models.Subscription, error) {
	receivedSubsc, _ := u.subscriptionRepository.GetByEmail(ctx, subscription.Email)

	if receivedSubsc != nil {
		u.logger.Warnf("Email %s already subscribed", subscription.Email)
		return nil, apperrors.AlreadySubscribedError
	}

	createdSubscription, err := u.subscriptionRepository.Create(ctx, subscription)
	if err != nil {
		return nil, err
	}

	return createdSubscription, nil
}

func (u *SubscriptionUseCase) Confirm(ctx context.Context, token string) (*models.Subscription, error) {
	receivedSubsc, err := u.subscriptionRepository.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if receivedSubsc == nil {
		return nil, apperrors.TokenNotFoundError
	}

	receivedSubsc.Confirmed = true
	updatedSubsc, err := u.subscriptionRepository.Update(ctx, *receivedSubsc)

	if err != nil {
		return nil, err
	}

	return updatedSubsc, nil

}

func (u *SubscriptionUseCase) Unsubscribe(ctx context.Context, token string) error {
	receivedSubsc, err := u.subscriptionRepository.GetByToken(ctx, token)
	if err != nil {
		return err
	}

	if receivedSubsc == nil {
		return apperrors.TokenNotFoundError
	}

	err = u.subscriptionRepository.DeleteByToken(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (u *SubscriptionUseCase) ListByFrequency(ctx context.Context, frequency models.Frequency, lastID, pageSize int) ([]models.Subscription, error) {
	receivedSubscriptions, err := u.subscriptionRepository.ListConfirmedByFrequency(ctx, frequency, lastID, pageSize)
	if err != nil {
		return nil, err
	}

	return receivedSubscriptions, nil
}
