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
	log := s.logger.WithContext(ctx)

	receivedSubsc, err := s.subscriptionRepository.GetByEmail(ctx, subscription.Email)
	if err != nil {
		return nil, err
	}
	if receivedSubsc != nil {
		log.Infof("Subscription attempt stopped: email %s already subscribed", subscription.Email)
		return nil, domainerrors.ErrAlreadySubscribed
	}

	log.Debugf("Generating confirmation token for email: %s", subscription.Email)
	token := s.tokenManager.Generate(ctx)
	subscription.Token = token

	createdSubscription, err := s.subscriptionRepository.Create(ctx, *subscription)
	if err != nil {
		return nil, err
	}
	log.Infof("Subscription created in database: id=%d, email=%s", createdSubscription.ID, createdSubscription.Email)

	confirmationInfo := contracts.ConfirmationInfo{
		Email:     createdSubscription.Email,
		Token:     createdSubscription.Token,
		Frequency: createdSubscription.Frequency,
	}

	log.Infof("Sending confirmation email: email=%s, token=%s", createdSubscription.Email, createdSubscription.Token)
	s.mailer.SendConfirmation(ctx, &confirmationInfo)

	return createdSubscription, nil
}

func (s *SubscriptionService) Confirm(ctx context.Context, token string) error {
	log := s.logger.WithContext(ctx)

	log.Debugf("Validating token for confirmation: %s", token)
	tokenIsValid := s.tokenManager.Validate(ctx, token)
	if !tokenIsValid {
		log.Warnf("Invalid token used for confirmation: %s", token)
		return domainerrors.ErrInvalidToken
	}

	receivedSubsc, err := s.subscriptionRepository.GetByToken(ctx, token)
	if err != nil {
		return err
	}
	if receivedSubsc == nil {
		log.Warnf("Token not found in database: %s", token)
		return domainerrors.ErrTokenNotFound
	}

	if !receivedSubsc.Confirmed {
		log.Infof("Confirming subscription: email=%s, token=%s", receivedSubsc.Email, token)
		receivedSubsc.Confirmed = true
		updatedSubsc, err := s.subscriptionRepository.Update(ctx, *receivedSubsc)
		if err != nil {
			return err
		}

		log.Infof("Subscription confirmed in database: id=%d, email=%s", updatedSubsc.ID, updatedSubsc.Email)

		confirmedInfo := contracts.ConfirmedInfo{
			Email:     updatedSubsc.Email,
			Token:     updatedSubsc.Token,
			Frequency: updatedSubsc.Frequency,
		}

		log.Infof("Sending confirmation success email: %s", updatedSubsc.Email)
		s.mailer.SendConfirmed(ctx, &confirmedInfo)
	}

	return nil
}

func (s *SubscriptionService) Unsubscribe(ctx context.Context, token string) error {
	log := s.logger.WithContext(ctx)

	log.Debugf("Validating token for unsubscription: %s", token)
	tokenIsValid := s.tokenManager.Validate(ctx, token)
	if !tokenIsValid {
		return domainerrors.ErrInvalidToken
	}

	receivedSubsc, err := s.subscriptionRepository.GetByToken(ctx, token)
	if err != nil {
		return err
	}
	if receivedSubsc == nil {
		log.Warnf("Token not found in database: %s", token)
		return domainerrors.ErrTokenNotFound
	}

	err = s.subscriptionRepository.DeleteByToken(ctx, token)
	if err != nil {
		return err
	}

	log.Infof("Subscription deleted from database: id=%d, email=%s", receivedSubsc.ID, receivedSubsc.Email)

	unsubscribeInfo := contracts.UnsubscribeInfo{
		Email:     receivedSubsc.Email,
		City:      receivedSubsc.City,
		Frequency: receivedSubsc.Frequency,
	}

	log.Infof("Sending unsubscription success email: %s", receivedSubsc.Email)
	s.mailer.SendUnsubscribed(ctx, &unsubscribeInfo)

	return nil
}

func (s *SubscriptionService) ListByFrequency(ctx context.Context, query *models.ListSubscriptionsQuery) ([]models.Subscription, error) {
	log := s.logger.WithContext(ctx)

	receivedSubscriptions, err := s.subscriptionRepository.ListConfirmedByFrequency(ctx, query.Frequency, query.LastID, query.PageSize)
	if err != nil {
		return nil, err
	}

	log.Infof("Subscription list received from database: frequency=%s", query.Frequency)

	return receivedSubscriptions, nil
}
