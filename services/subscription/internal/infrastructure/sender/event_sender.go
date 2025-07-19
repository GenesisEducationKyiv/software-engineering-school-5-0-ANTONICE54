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
	e := mappers.ConfirmationInfoToEvent(info)
	err := s.publisher.Publish(ctx, e)
	if err != nil {
		s.logger.Warnf("failed to publish event: %s", err.Error())
	}
}
func (s *EventSender) SendConfirmed(ctx context.Context, info *contracts.ConfirmedInfo) {
	e := mappers.ConfirmedInfoToEvent(info)
	err := s.publisher.Publish(ctx, e)
	if err != nil {
		s.logger.Warnf("failed to publish event: %s", err.Error())
	}
}

func (s *EventSender) SendUnsubscribed(ctx context.Context, info *contracts.UnsubscribeInfo) {
	e := mappers.UnsubscribedInfoToEvent(info)
	err := s.publisher.Publish(ctx, e)
	if err != nil {
		s.logger.Warnf("failed to publish event: %s", err.Error())
	}
}
