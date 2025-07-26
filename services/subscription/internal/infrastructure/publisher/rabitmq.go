package publisher

import (
	"context"
	infraerror "subscription-service/internal/infrastructure/errors"
	"subscription-service/internal/infrastructure/events"
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

	routeKey, err := event.RoutingKey()

	if err != nil {
		return err

	}

	err = p.ch.PublishWithContext(
		ctx,
		p.exchange,
		routeKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/x-protobuf",
			Body:        event.Body,
		})

	if err != nil {
		return infraerror.ErrInternal
	}

	return nil
}
