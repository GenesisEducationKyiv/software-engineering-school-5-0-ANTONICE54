package main

import (
	"log"
	"subscription-service/internal/config"
	"subscription-service/internal/domain/usecases"
	"subscription-service/internal/infrastructure/database"
	"subscription-service/internal/infrastructure/decorators"
	"subscription-service/internal/infrastructure/metrics"
	"subscription-service/internal/infrastructure/repositories"
	"subscription-service/internal/infrastructure/sender"

	"subscription-service/internal/infrastructure/token"
	"subscription-service/internal/presentation/server"
	"subscription-service/internal/presentation/server/handlers"
	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/publisher"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to read from config: %s", err.Error())
	}
	logrusLog := logger.NewLogrus(cfg.ServiceName)
	prometheusMetrics := metrics.NewPrometheus(logrusLog)

	conn, err := amqp.Dial(cfg.RabbitMQSource)
	if err != nil {
		logrusLog.Fatalf("Failed to connect to RabbitMQ: %s", err.Error())
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logrusLog.Errorf("Failed to close RabbitMQ connection: %v", err)
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}
	defer func() {
		if err := ch.Close(); err != nil {
			logrusLog.Errorf("Failed to close RabbitMQ channel: %v", err)
		}
	}()

	db := database.Connect(cfg)
	database.RunMigration(db)

	subscRepo := repositories.NewSubscriptionRepository(db, logrusLog)
	tokenManager := token.NewUUIDManager()

	rabbitMQPublisher := publisher.NewRabbitMQPublisher(ch, cfg.Exchange, logrusLog)

	eventSender := sender.NewEventSender(rabbitMQPublisher, logrusLog)

	subscUseCase := usecases.NewSubscriptionService(subscRepo, tokenManager, eventSender, logrusLog)
	metricSubscUseCase := decorators.NewSubscriptionServiceMetricsDecorator(*subscUseCase, prometheusMetrics, logrusLog)
	subscHandler := handlers.NewSubscriptionHandler(metricSubscUseCase, logrusLog)

	go prometheusMetrics.StartMetricsServer(cfg.MetricsServerPort)

	app := server.New(subscHandler, logrusLog)

	if err := app.Start(cfg.GRPCPort); err != nil {
		logrusLog.Fatalf("Failed to start server: %v", err)
	}
}
