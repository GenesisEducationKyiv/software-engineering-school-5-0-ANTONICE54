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
	event, err := events.NewConfirmation(info)
	if err != nil {
		s.logger.Warnf("failed to create event: %s", err.Error())
		return
	}

	routingKey, err := event.RoutingKey()
	if err != nil {
		s.logger.Warnf("failed to get routing key: %s", err.Error())
	}

	err = s.publisher.Publish(ctx, routingKey, event.Body)
	if err != nil {
		s.logger.Warnf("failed to publish event: %s", err.Error())
	}
	s.logger.Infof("Confirmation event published")
}

func (s *EventSender) SendConfirmed(ctx context.Context, info *contracts.ConfirmedInfo) {
	event, err := events.NewConfirmed(info)
	if err != nil {
		s.logger.Warnf("failed to create event: %s", err.Error())
		return
	}

	routingKey, err := event.RoutingKey()
	if err != nil {
		s.logger.Warnf("failed to get routing key: %s", err.Error())
	}

	err = s.publisher.Publish(ctx, routingKey, event.Body)
	if err != nil {
		s.logger.Warnf("failed to publish event: %s", err.Error())
	}
	s.logger.Infof("Confirmed event published")

}

func (s *EventSender) SendUnsubscribed(ctx context.Context, info *contracts.UnsubscribeInfo) {
	event, err := events.NewUnsubscribed(info)
	if err != nil {
		s.logger.Warnf("failed to create event: %s", err.Error())
		return
	}

	routingKey, err := event.RoutingKey()
	if err != nil {
		s.logger.Warnf("failed to get routing key: %s", err.Error())
	}

	err = s.publisher.Publish(ctx, routingKey, event.Body)
	if err != nil {
		s.logger.Warnf("failed to publish event: %s", err.Error())
	}
}
