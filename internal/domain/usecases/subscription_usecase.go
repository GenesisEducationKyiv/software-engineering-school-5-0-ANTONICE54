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
	}

	TokenManager interface {
		Generate(ctx context.Context) string
		Validate(ctx context.Context, token string) bool
	}

	NotificationService interface {
		SendConfirmation(ctx context.Context, email, token string, frequency models.Frequency)
		SendConfirmed(ctx context.Context, email, token string, frequency models.Frequency)
	}

	SubscriptionUseCase struct {
		subscriptionRepository SubscriptionRepository
		tokenManager           TokenManager
		mailer                 NotificationService
		logger                 logger.Logger
	}
)

func NewSubscriptionUseCase(subscriptionRepo SubscriptionRepository, tokenManager TokenManager, mailer NotificationService, logger logger.Logger) *SubscriptionUseCase {
	return &SubscriptionUseCase{
		subscriptionRepository: subscriptionRepo,
		tokenManager:           tokenManager,
		mailer:                 mailer,
		logger:                 logger,
	}
}

func (uc *SubscriptionUseCase) Subscribe(ctx context.Context, email, frequency, city string) (*models.Subscription, error) {

	receivedSubsc, _ := uc.subscriptionRepository.GetByEmail(ctx, email)

	if receivedSubsc != nil {
		uc.logger.Warnf("Email %s already subscribed", email)
		return nil, apperrors.AlreadySubscribedError
	}

	token := uc.tokenManager.Generate(ctx)
	subscription := models.Subscription{
		Email:     email,
		Frequency: models.Frequency(frequency),
		City:      city,
		Confirmed: false,
		Token:     token,
	}

	createdSubscription, err := uc.subscriptionRepository.Create(ctx, subscription)
	if err != nil {
		return nil, err
	}

	uc.mailer.SendConfirmation(ctx, createdSubscription.Email, createdSubscription.Token, createdSubscription.Frequency)

	return createdSubscription, nil
}

func (uc *SubscriptionUseCase) Confirm(ctx context.Context, token string) error {
	tokenIsValid := uc.tokenManager.Validate(ctx, token)
	if !tokenIsValid {
		return apperrors.InvalidTokenError
	}

	receivedSubsc, err := uc.subscriptionRepository.GetByToken(ctx, token)
	if err != nil {
		return err
	}
	if receivedSubsc == nil {
		return apperrors.TokenNotFoundError
	}

	if !receivedSubsc.Confirmed {
		receivedSubsc.Confirmed = true
		updatedSubsc, err := uc.subscriptionRepository.Update(ctx, *receivedSubsc)

		if err != nil {
			return err
		}

		uc.mailer.SendConfirmed(ctx, updatedSubsc.Email, updatedSubsc.Token, updatedSubsc.Frequency)
	}

	return nil
}

func (uc *SubscriptionUseCase) Unsubscribe(ctx context.Context, token string) error {
	tokenIsValid := uc.tokenManager.Validate(ctx, token)
	if !tokenIsValid {
		return apperrors.InvalidTokenError
	}

	receivedSubsc, err := uc.subscriptionRepository.GetByToken(ctx, token)
	if err != nil {
		return err
	}
	if receivedSubsc == nil {
		return apperrors.TokenNotFoundError
	}

	err = uc.subscriptionRepository.DeleteByToken(ctx, token)
	if err != nil {
		return err
	}

	return nil
}
