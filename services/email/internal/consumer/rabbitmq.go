package consumer

import (
	"context"
	"sync"
	"weather-forecast/pkg/ctxutil"
	"weather-forecast/pkg/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

const maxWorkers = 10

type (
	EventHandler interface {
		Handle(ctx context.Context, routingKey string, body []byte)
	}

	Consumer struct {
		conn         *amqp.Connection
		exchangeName string
		handler      EventHandler
		channel      *amqp.Channel
		wg           *sync.WaitGroup
		cancel       context.CancelFunc
		semaphore    chan struct{}
		logger       logger.Logger
	}
)

func NewConsumer(conn *amqp.Connection, exchange string, handler EventHandler, logger logger.Logger) *Consumer {
	return &Consumer{
		conn:         conn,
		exchangeName: exchange,
		handler:      handler,
		wg:           &sync.WaitGroup{},
		semaphore:    make(chan struct{}, maxWorkers),
		logger:       logger,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	c.logger.Infof("Starting RabbitMQ consumer for exchange: %s", c.exchangeName)

	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	c.channel = ch

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

		for {
			select {

			case <-ctx.Done():
				c.logger.Infof("Shutting down consumer goroutine...")
				return

			case d, ok := <-msgs:
				if !ok {
					c.logger.Infof("Message channel closed")
					return
				}

				c.semaphore <- struct{}{}

				c.wg.Add(1)
				go func(d amqp.Delivery) {
					defer func() {
						<-c.semaphore
						c.wg.Done()
					}()
					c.handleMessage(d)
				}(d)
			}
		}
	}()

	return nil
}

func (c *Consumer) handleMessage(d amqp.Delivery) {

	c.logger.Infof("Event received: %s", d.RoutingKey)

	msgCtx := context.Background()
	if val, ok := d.Headers[ctxutil.CorrelationIDKey.String()]; ok {
		if correlationID, ok := val.(string); ok {
			//nolint:staticcheck
			msgCtx = context.WithValue(context.Background(), ctxutil.CorrelationIDKey.String(), correlationID)
		}
	} else {
		c.logger.Warnf("correlation-id not found in headers")
	}

	c.handler.Handle(msgCtx, d.RoutingKey, d.Body)
}

func (c *Consumer) Stop() {
	c.logger.Infof("Closing RabbitMQ channel...")
	if c.cancel != nil {
		c.cancel()
	}

	c.wg.Wait()

	if c.channel != nil {
		err := c.channel.Close()
		if err != nil {
			c.logger.Errorf("Failed to close RabbitMQ channel.")
		}
	}

	c.logger.Infof("Consumer stopped")
}
