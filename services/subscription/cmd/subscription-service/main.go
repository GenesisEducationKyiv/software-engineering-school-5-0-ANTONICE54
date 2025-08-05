package main

import (
	"log"
	"os"
	"os/signal"
	"subscription-service/internal/config"
	"subscription-service/internal/domain/usecases"
	"subscription-service/internal/infrastructure/database"
	"subscription-service/internal/infrastructure/decorators"
	"subscription-service/internal/infrastructure/metrics"
	"subscription-service/internal/infrastructure/repositories"
	"subscription-service/internal/infrastructure/sender"
	"syscall"

	"subscription-service/internal/infrastructure/token"
	"subscription-service/internal/presentation/server"
	"subscription-service/internal/presentation/server/handlers"
	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/publisher"
	"weather-forecast/pkg/rabbitmq"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to read from config: %s", err.Error())
	}
	logSampler := logger.NewRateSampler(cfg.LogSamplingRate)
	logrusLog, err := logger.NewLogrus(cfg.ServiceName, cfg.LogLevel, logSampler)
	if err != nil {
		log.Fatalf("Failed to initialize logger with level '%s': %v", cfg.LogLevel, err)
	}
	prometheusMetrics := metrics.NewPrometheus(logrusLog)

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

	db, err := database.Connect(&cfg.DB)
	if err != nil {
		logrusLog.Fatalf("Failed to establish connection with database: %s", err.Error())
	}

	err = database.RunMigration(db)
	if err != nil {
		logrusLog.Fatalf("Failed to migrate database: %s", err.Error())
	}

	subscRepo := repositories.NewSubscriptionRepository(db, logrusLog)
	tokenManager := token.NewUUIDManager()

	rabbitMQPublisher := publisher.NewRabbitMQPublisher(ch, cfg.RabbitMQ.Exchange, logrusLog)

	eventSender := sender.NewEventSender(rabbitMQPublisher, logrusLog)

	subscUseCase := usecases.NewSubscriptionService(subscRepo, tokenManager, eventSender, logrusLog)
	metricSubscUseCase := decorators.NewSubscriptionServiceMetricsDecorator(*subscUseCase, prometheusMetrics, logrusLog)
	subscHandler := handlers.NewSubscriptionHandler(metricSubscUseCase, logrusLog)

	go prometheusMetrics.StartMetricsServer(cfg.MetricsServerPort)

	app := server.New(subscHandler, logrusLog)
	go func() {
		if err := app.Start(cfg.GRPCPort); err != nil {
			logrusLog.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrusLog.Infof("Shutting down subscription service...")
	app.Shutdown()
	logrusLog.Infof("Service stopped gracefully")

}
