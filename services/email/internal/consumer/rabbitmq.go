package consumer

import (
	"context"
	"weather-forecast/pkg/ctxutil"
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

	c.logger.Infof("Starting RabbitMQ consumer for exchange: %s", c.exchangeName)

	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	c.logger.Debugf("Declaring exchange: %s", c.exchangeName)
	err = ch.ExchangeDeclare(c.exchangeName, "topic", true, false, false, false, nil)
	if err != nil {
		return err
	}

	c.logger.Debugf("Declaring queue")
	q, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return err
	}
	c.logger.Debugf("Created queue: %s", q.Name)

	c.logger.Debugf("Binding queue %s to exchange %s with routing key '#'", q.Name, c.exchangeName)
	err = ch.QueueBind(q.Name, "#", c.exchangeName, false, nil)
	if err != nil {
		return err
	}

	c.logger.Debugf("Starting to consume messages from queue: %s", q.Name)
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		c.logger.Infof("RabbitMQ consumer started successfully, waiting for messages...")

		for d := range msgs {
			c.logger.Infof("Event received: %s", d.RoutingKey)

			processID := "unknown-process"
			if val, ok := d.Headers[ctxutil.ProcessIDKey.String()]; ok {
				if s, ok := val.(string); ok {
					processID = s
				}
			}
			ctx = context.WithValue(context.Background(), ctxutil.ProcessIDKey.String(), processID)

			c.handler.Handle(ctx, d.RoutingKey, d.Body)
		}
	}()

	return nil
}
