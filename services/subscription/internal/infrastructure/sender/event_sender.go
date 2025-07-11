package sender

import (
	"context"
	"subscription-service/internal/domain/dto"
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

func (s *EventSender) SendConfirmation(ctx context.Context, info *dto.ConfirmationInfo) {
	e := mappers.ConfirmationInfoToEvent(info)
	err := s.publisher.Publish(ctx, e)
	if err != nil {
		s.logger.Warnf("failed to publish event: %s", err.Error())
	}
}
func (s *EventSender) SendConfirmed(ctx context.Context, info *dto.ConfirmedInfo) {
	e := mappers.ConfirmedInfoToEvent(info)
	err := s.publisher.Publish(ctx, e)
	if err != nil {
		s.logger.Warnf("failed to publish event: %s", err.Error())
	}
}
