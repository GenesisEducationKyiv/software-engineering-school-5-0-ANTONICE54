package consumer

import (
	"context"
	"weather-forecast/pkg/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	EventHandler interface {
		Handle(ctx context.Context, routingKey string, body []byte)
	}

	Consumer struct {
		conn         *amqp.Connection
		exchangeName string
		handler      EventHandler
		logger       logger.Logger
	}
)

func NewConsumer(conn *amqp.Connection, exchange string, handler EventHandler, logger logger.Logger) *Consumer {
	return &Consumer{
		conn:         conn,
		exchangeName: exchange,
		handler:      handler,
		logger:       logger,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(c.exchangeName, "topic", true, false, false, false, nil)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return err
	}

	err = ch.QueueBind(q.Name, "#", c.exchangeName, false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			c.logger.Infof("Event received:", d.RoutingKey)
			c.handler.Handle(ctx, d.RoutingKey, d.Body)
		}
	}()

	return nil
}
