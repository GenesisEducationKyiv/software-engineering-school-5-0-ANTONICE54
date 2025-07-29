package sender

import (
	"context"
	"subscription-service/internal/domain/contracts"
	"subscription-service/internal/infrastructure/mappers"
	"weather-forecast/pkg/events"
	"weather-forecast/pkg/logger"
)

type (
	EventPublisher interface {
		Publish(ctx context.Context, event events.Event) error
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

	log.Infof("Publishing confirmation event: email=%s, token=%s", info.Email, info.Token)
	e := mappers.ConfirmationInfoToEvent(info)
	err := s.publisher.Publish(ctx, e)
	if err != nil {
		log.Errorf("Failed to publish confirmation event for email %s: %v", info.Email, err)
	}

	log.Debugf("Confirmation event published successfully: email=%s", info.Email)
}
func (s *EventSender) SendConfirmed(ctx context.Context, info *contracts.ConfirmedInfo) {
	log := s.logger.WithContext(ctx)

	log.Infof("Publishing confirmed event: email=%s, token=%s", info.Email, info.Token)

	e := mappers.ConfirmedInfoToEvent(info)
	err := s.publisher.Publish(ctx, e)
	if err != nil {
		log.Errorf("Failed to publish confirmed event for email %s: %v", info.Email, err)
	}

	log.Debugf("Confirmed event published successfully: email=%s", info.Email)
}

func (s *EventSender) SendUnsubscribed(ctx context.Context, info *contracts.UnsubscribeInfo) {
	log := s.logger.WithContext(ctx)

	log.Infof("Publishing unsubscribed event: email=%s, city=%s", info.Email, info.City)

	e := mappers.UnsubscribedInfoToEvent(info)
	err := s.publisher.Publish(ctx, e)
	if err != nil {
		log.Errorf("Failed to publish unsubscribed event for email %s: %v", info.Email, err)
	}

	log.Debugf("Unsubscribed event published successfully: email=%s", info.Email)
}
