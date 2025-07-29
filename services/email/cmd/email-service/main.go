package main

import (
	"context"
	"email-service/internal/config"
	"email-service/internal/consumer"
	"email-service/internal/mailer"
	"email-service/internal/mailer/decorators"
	"email-service/internal/metrics"
	"email-service/internal/processors"
	"email-service/internal/services"
	"log"
	"os"
	"os/signal"
	"syscall"
	"weather-forecast/pkg/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to read from config: %s", err.Error())
	}
	logrusLog := logger.NewLogrus(cfg.ServiceName)
	prometheusMetrics := metrics.NewPrometheus(logrusLog)

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
	retryMailer := decorators.NewRetryMailer(mailer, cfg.Retry, logrusLog)
	metricsMailer := decorators.NewMetricMailer(retryMailer, prometheusMetrics, logrusLog)
	emailBuilder := services.NewSimpleEmailBuild(cfg.ServerHost, logrusLog)
	notificationService := services.NewNotificationService(metricsMailer, emailBuilder, logrusLog)

	eventProcessor := processors.NewEventProcessor(notificationService, logrusLog)

	rabbitMQConsumer := consumer.NewConsumer(conn, cfg.Exchange, eventProcessor, logrusLog)

	go prometheusMetrics.StartMetricsServer(cfg.MetricsServerPort)

	if err := rabbitMQConsumer.Start(context.Background()); err != nil {
		logrusLog.Fatalf("Failed to start consumer: %v", err)
	}

	logrusLog.Infof("Email service started successfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrusLog.Infof("Shutting down email service...")

}
