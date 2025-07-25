package rabbitmq

import (
	"fmt"
	"time"
	"weather-forecast/pkg/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectWithRetry(cfg Config, log logger.Logger) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	for i := 0; i < cfg.Retries; i++ {
		log.Infof("Attempting to connect to RabbitMQ (attempt %d/%d)...", i+1, cfg.Retries)

		conn, err = amqp.Dial(cfg.Source)
		if err == nil {
			log.Infof("Successfully connected to RabbitMQ")
			return conn, nil
		}

		log.Warnf("Failed to connect to RabbitMQ: %v", err)

		if i < cfg.Retries-1 {
			log.Infof("Retrying in %v...", cfg.RetryDelay)
			time.Sleep(time.Duration(cfg.RetryDelay) * time.Second)
		}
	}

	return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %s", cfg.Retries, err.Error())
}
