package main

import (
	"log"
	"subscription-service/internal/config"
	"subscription-service/internal/domain/usecases"
	"subscription-service/internal/infrastructure/database"
	"subscription-service/internal/infrastructure/publisher"
	"subscription-service/internal/infrastructure/repositories"
	"subscription-service/internal/infrastructure/sender"
	"subscription-service/internal/infrastructure/token"
	"subscription-service/internal/presentation/server"
	"subscription-service/internal/presentation/server/handlers"
	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/rabbitmq"
)

func main() {

	logrusLog := logger.NewLogrus()

	cfg, err := config.Load(logrusLog)
	if err != nil {
		logrusLog.Fatalf("Failed to read from config: %s", err.Error())
	}

	conn, err := rabbitmq.ConnectWithRetry(cfg.RabbitMQ, logrusLog)
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

	db := database.Connect(&cfg.DB)
	database.RunMigration(db)

	subscRepo := repositories.NewSubscriptionRepository(db, logrusLog)
	tokenManager := token.NewUUIDManager()

	rabbitMQPublisher := publisher.NewRabbitMQPublisher(ch, cfg.RabbitMQ.Exchange, logrusLog)

	eventSender := sender.NewEventSender(rabbitMQPublisher, logrusLog)

	subscUseCase := usecases.NewSubscriptionService(subscRepo, tokenManager, eventSender, logrusLog)
	subscHandler := handlers.NewSubscriptionHandler(subscUseCase, logrusLog)

	app := server.New(subscHandler, logrusLog)

	if err := app.Start(cfg.GRPCPort); err != nil {
		logrusLog.Fatalf("Failed to start server: %v", err)
	}
}
