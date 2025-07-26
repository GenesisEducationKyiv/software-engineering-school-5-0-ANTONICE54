package sender

import (
	"context"
	"subscription-service/internal/domain/contracts"
	"subscription-service/internal/infrastructure/events"

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
	event, err := events.NewConfirmation(info)

	if err != nil {
		s.logger.Warnf("failed to create event: %s", err.Error())
		return
	}
	err = s.publisher.Publish(ctx, *event)
	if err != nil {
		s.logger.Warnf("failed to publish event: %s", err.Error())
	}
}

func (s *EventSender) SendConfirmed(ctx context.Context, info *contracts.ConfirmedInfo) {
	event, err := events.NewConfirmed(info)
	if err != nil {
		s.logger.Warnf("failed to create event: %s", err.Error())
		return
	}

	err = s.publisher.Publish(ctx, *event)
	if err != nil {
		s.logger.Warnf("failed to publish event: %s", err.Error())
	}
}
