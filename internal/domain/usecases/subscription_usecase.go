package usecases

import (
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/apperrors"
	"weather-forecast/internal/infrastructure/logger"
)

type (
	SubscriptionRepositoryI interface {
		Create(subscription models.Subscription) (*models.Subscription, error)
		GetByEmail(email string) (*models.Subscription, error)
		GetByToken(token string) (*models.Subscription, error)
		Update(subscription models.Subscription) (*models.Subscription, error)
		DeleteByToken(token string) error
		ListConfirmedByFrequency(frequency models.Frequency) ([]models.Subscription, error)
	}
	SubscriptionUseCase struct {
		subscriptionRepository SubscriptionRepositoryI
		logger                 logger.Logger
	}
)

func NewSubscriptionUseCase(subscRepo SubscriptionRepositoryI, logger logger.Logger) *SubscriptionUseCase {
	return &SubscriptionUseCase{
		subscriptionRepository: subscRepo,
		logger:                 logger,
	}
}

func (u *SubscriptionUseCase) Subscribe(subscription models.Subscription) (*models.Subscription, error) {
	receivedSubsc, _ := u.subscriptionRepository.GetByEmail(subscription.Email)

	if receivedSubsc != nil {
		u.logger.Warnf("Email %s already subscribed", subscription.Email)
		return nil, apperrors.AlreadySubscribedError
	}

	createdSubscription, err := u.subscriptionRepository.Create(subscription)
	if err != nil {
		return nil, err
	}

	return createdSubscription, nil
}

func (u *SubscriptionUseCase) Confirm(token string) (*models.Subscription, error) {
	receivedSubsc, err := u.subscriptionRepository.GetByToken(token)
	if err != nil {
		return nil, err
	}

	if receivedSubsc == nil {
		return nil, apperrors.TokenNotFoundError
	}

	receivedSubsc.Confirmed = true
	updatedSubsc, err := u.subscriptionRepository.Update(*receivedSubsc)

	if err != nil {
		return nil, err
	}

	return updatedSubsc, nil

}

func (u *SubscriptionUseCase) Unsubscribe(token string) error {
	receivedSubsc, err := u.subscriptionRepository.GetByToken(token)
	if err != nil {
		return err
	}

	if receivedSubsc == nil {
		return apperrors.TokenNotFoundError
	}

	err = u.subscriptionRepository.DeleteByToken(token)
	if err != nil {
		return err
	}

	return nil
}

func (u *SubscriptionUseCase) ListByFrequency(frequency models.Frequency) ([]models.Subscription, error) {
	receivedSubscriptions, err := u.subscriptionRepository.ListConfirmedByFrequency(frequency)
	if err != nil {
		return nil, err
	}

	return receivedSubscriptions, nil
}
