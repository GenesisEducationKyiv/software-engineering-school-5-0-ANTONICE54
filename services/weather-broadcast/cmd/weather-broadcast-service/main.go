package main

import (
	"context"
	"log"
	"time"
	"weather-broadcast-service/internal/clients"
	"weather-broadcast-service/internal/publisher"
	"weather-broadcast-service/internal/scheduler"
	"weather-broadcast-service/internal/sender"
	"weather-broadcast-service/internal/services"
	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/proto/subscription"
	"weather-forecast/pkg/proto/weather"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	logrusLog := logger.NewLogrus()

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		logrusLog.Fatalf("Failed to read from config: %s", err.Error())
	}

	weatherConn, err := grpc.NewClient("weather-service:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logrusLog.Fatalf("Failed to connect to Subscription Service: %v", err)
	}
	defer weatherConn.Close()

	weatherGRPCClient := weather.NewWeatherServiceClient(weatherConn)
	weatherClient := clients.NewWeatherGRPCClient(weatherGRPCClient, logrusLog)

	subscConn, err := grpc.NewClient("subscription-service:8082",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logrusLog.Fatalf("Failed to connect to Subscription Service: %v", err)
	}
	defer subscConn.Close()

	subscriptionGRPCClient := subscription.NewSubscriptionServiceClient(subscConn)
	sunbscriptionClient := clients.NewSubscriptionGRPCClient(subscriptionGRPCClient, logrusLog)

	rabitMQSource := viper.GetString("RABBIT_MQ_SOURCE")
	conn, err := amqp.Dial(rabitMQSource)
	if err != nil {
		logrusLog.Fatalf("Failed to connect to RabbitMQ: %s", err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}
	defer ch.Close()

	exchange := viper.GetString("EXCHANGE")
	rabbitMQPublisher := publisher.NewRabbitMQPublisher(ch, exchange, logrusLog)
	eventSender := sender.NewEventSender(rabbitMQPublisher, logrusLog)

	timezone := viper.GetString("TIMEZONE")
	location, err := time.LoadLocation(timezone)
	if err != nil {
		logrusLog.Fatalf("Failed to load timezone: %s", err.Error())
	}

	weatherBroadcastService := services.NewWeatherBroadcastService(sunbscriptionClient, weatherClient, eventSender, logrusLog)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler := scheduler.New(ctx, weatherBroadcastService, location, logrusLog)
	scheduler.SetUp()
	scheduler.Run()

	select {}
}
