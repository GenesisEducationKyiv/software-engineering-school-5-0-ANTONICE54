package main

import (
	"context"
	"email-service/internal/config"
	"email-service/internal/consumer"
	"email-service/internal/mailer"
	"email-service/internal/processors"
	"email-service/internal/services"
	"os"
	"os/signal"
	"syscall"
	"weather-forecast/pkg/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	logrusLog := logger.NewLogrus()

	cfg, err := config.Load(logrusLog)
	if err != nil {
		logrusLog.Fatalf("Failed to read from config: %s", err.Error())
	}

	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		logrusLog.Fatalf("Failed to connect to RabbitMQ: %s", err.Error())
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logrusLog.Errorf("Failed to close RabbitMQ connection: %v", err)
		}
	}()

	mailer := mailer.NewSMTPMailer(cfg, logrusLog)
	notificationService := services.NewNotificationService(mailer, cfg.ServerHost, logrusLog)

	eventProcessor := processors.NewEventProcessor(notificationService, logrusLog)

	rabbitMQConsumer := consumer.NewConsumer(conn, cfg.Exchange, eventProcessor, logrusLog)

	if err := rabbitMQConsumer.Start(context.Background()); err != nil {
		logrusLog.Fatalf("Failed to start consumer: %v", err)
	}

	logrusLog.Info("Email service started successfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrusLog.Info("Shutting down email service...")

}
