package sender

import (
	"context"
	"subscription-service/internal/domain/contracts"
	"subscription-service/internal/infrastructure/events"

	"weather-forecast/pkg/logger"
)

type (
	EventPublisher interface {
		Publish(ctx context.Context, routingKey string, body []byte) error
	}

	EventSender struct {
		publisher EventPublisher
		logger    logger.Logger
	}
)

func NewEventSender(publisher EventPublisher, logger logger.Logger) *EventSender {
	return &EventSender{
		publisher: publisher,
		logger:    logger,
	}
}

func (s *EventSender) SendConfirmation(ctx context.Context, info *contracts.ConfirmationInfo) {
	log := s.logger.WithContext(ctx)

	log.Debugf("Creating confirmation event: email=%s, token=%s", info.Email, info.Token)
	event, err := events.NewConfirmation(info)
	if err != nil {
		log.Errorf("Failed to create confirmation event for email %s: %v", info.Email, err)

		return
	}

	routingKey, err := event.RoutingKey()
	if err != nil {
		log.Errorf("Failed to get confirmation event routing key for email %s: %v", info.Email, err)

		return
	}

	log.Infof("Publishing confirmation event: email=%s, token=%s", info.Email, info.Token)
	err = s.publisher.Publish(ctx, routingKey, event.Body)
	if err != nil {
		log.Errorf("Failed to publish confirmation event for email %s: %v", info.Email, err)
		return
	}

	log.Debugf("Confirmation event published successfully: email=%s", info.Email)
}

func (s *EventSender) SendConfirmed(ctx context.Context, info *contracts.ConfirmedInfo) {
	log := s.logger.WithContext(ctx)

	log.Debugf("Creating confirmed event: email=%s, token=%s", info.Email, info.Token)
	event, err := events.NewConfirmed(info)
	if err != nil {
		log.Errorf("Failed to create confirmed event for email %s: %v", info.Email, err)
		return
	}

	routingKey, err := event.RoutingKey()
	if err != nil {
		log.Errorf("Failed to get confirmed event routing key for email %s: %v", info.Email, err)
		return
	}

	log.Infof("Publishing confirmed event: email=%s, token=%s", info.Email, info.Token)
	err = s.publisher.Publish(ctx, routingKey, event.Body)
	if err != nil {
		log.Errorf("Failed to publish confirmed event for email %s: %v", info.Email, err)
		return
	}

	log.Debugf("Confirmed event published successfully: email=%s", info.Email)
}

func (s *EventSender) SendUnsubscribed(ctx context.Context, info *contracts.UnsubscribeInfo) {
	log := s.logger.WithContext(ctx)

	log.Debugf("Creating unsubscribed event: email=%s", info.Email)
	event, err := events.NewUnsubscribed(info)
	if err != nil {
		log.Errorf("Failed to create unsubscribed event for email %s: %v", info.Email, err)
		return
	}

	routingKey, err := event.RoutingKey()
	if err != nil {
		log.Errorf("Failed to get unsubscribed event routing key for email %s: %v", info.Email, err)
		return
	}

	log.Infof("Publishing unsubscribed event: email=%s", info.Email)
	err = s.publisher.Publish(ctx, routingKey, event.Body)
	if err != nil {
		log.Errorf("Failed to publish unsubscribed event for email %s: %v", info.Email, err)
		return
	}

	log.Debugf("Unsubscribed event published successfully: email=%s", info.Email)
}
