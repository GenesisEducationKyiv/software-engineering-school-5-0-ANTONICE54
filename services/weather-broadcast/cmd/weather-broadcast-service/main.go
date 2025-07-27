package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"weather-broadcast-service/internal/clients"
	"weather-broadcast-service/internal/config"
	"weather-broadcast-service/internal/scheduler"
	"weather-broadcast-service/internal/sender"
	"weather-broadcast-service/internal/services"
	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/proto/subscription"
	"weather-forecast/pkg/publisher"

	"weather-forecast/pkg/proto/weather"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to read from config: %s", err.Error())
	}
	logrusLog := logger.NewLogrus(cfg.ServiceName)

	weatherConn, err := grpc.NewClient(cfg.WeatherServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logrusLog.Fatalf("Failed to connect to Subscription Service: %v", err)
	}
	defer func() {
		if err := weatherConn.Close(); err != nil {
			logrusLog.Errorf("Failed to close gRPC connection with weather service: %v", err)
		}
	}()
	weatherGRPCClient := weather.NewWeatherServiceClient(weatherConn)
	weatherClient := clients.NewWeatherGRPCClient(weatherGRPCClient, logrusLog)

	subscConn, err := grpc.NewClient(cfg.SubscriptionServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logrusLog.Fatalf("Failed to connect to Subscription Service: %v", err)
	}
	defer func() {
		if err := subscConn.Close(); err != nil {
			logrusLog.Errorf("Failed to close gRPC connection with subscription service: %v", err)
		}
	}()

	subscriptionGRPCClient := subscription.NewSubscriptionServiceClient(subscConn)
	sunbscriptionClient := clients.NewSubscriptionGRPCClient(subscriptionGRPCClient, logrusLog)

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
	rabbitMQPublisher := publisher.NewRabbitMQPublisher(ch, cfg.Exchange, logrusLog)
	eventSender := sender.NewEventSender(rabbitMQPublisher, logrusLog)

	location, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		logrusLog.Fatalf("Failed to load timezone: %s", err.Error())
	}

	weatherBroadcastService := services.NewWeatherBroadcastService(sunbscriptionClient, weatherClient, eventSender, logrusLog)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler := scheduler.New(ctx, weatherBroadcastService, location, logrusLog)
	scheduler.SetUp()
	scheduler.Run()

	logrusLog.Infof("Weather broadcast service started successfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrusLog.Infof("Shutting down weather broadcast service...")
}
