package publisher

import (
	"context"
	"weather-forecast/pkg/ctxutil"
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

func (p *RabbitMQPublisher) Publish(ctx context.Context, routingKey string, body []byte) error {
	processID := ctxutil.GetProcessID(ctx)

	return p.ch.PublishWithContext(
		ctx,
		p.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/x-protobuf",
			Body:        body,
			Headers: amqp.Table{
				ctxutil.ProcessIDKey.String(): processID,
			},
		})

}
