package publisher

import (
	"context"
	"encoding/json"
	"weather-broadcast-service/internal/errors"
	"weather-forecast/pkg/events"
	"weather-forecast/pkg/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	RabbitMQPublisher struct {
		ch       *amqp.Channel
		exchange string
		logger   logger.Logger
	}
)

func NewRabbitMQPublisher(ch *amqp.Channel, exchange string, logger logger.Logger) *RabbitMQPublisher {
	return &RabbitMQPublisher{
		ch:       ch,
		exchange: exchange,
		logger:   logger,
	}
}

func (p *RabbitMQPublisher) Publish(ctx context.Context, event events.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		p.logger.Warnf("failed to marshal event: %w", err)
		return errors.InternalError
	}

	return p.ch.PublishWithContext(
		ctx,
		p.exchange,
		string(event.EventType()),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
