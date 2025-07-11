package main

import (
	"context"
	"email-service/internal/consumer"
	"email-service/internal/handlers"
	"email-service/internal/mailer"
	"email-service/internal/services"
	"weather-forecast/pkg/logger"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

func main() {
	logrusLog := logger.NewLogrus()

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		logrusLog.Fatalf("Failed to read from config: %s", err.Error())
	}

	//TODO: add rabbit source to env
	rabitMQSource := viper.GetString("RABBIT_MQ_SOURCE")
	conn, err := amqp.Dial(rabitMQSource)
	if err != nil {
		logrusLog.Fatalf("Failed to connect to RabbitMQ: %s", err.Error())
	}
	defer conn.Close()

	serverHost := viper.GetString("SERVER_HOST")
	mailerFrom := viper.GetString("MAILER_FROM")
	mailerHost := viper.GetString("MAILER_HOST")
	mailerPort := viper.GetString("MAILER_PORT")
	mailerUsername := viper.GetString("MAILER_USERNAME")
	mailerPassword := viper.GetString("MAILER_PASSWORD")
	mailer := mailer.NewSMTPMailer(mailerFrom, mailerHost, mailerPort, mailerUsername, mailerPassword, logrusLog)
	notificationService := services.NewNotificationService(mailer, serverHost, logrusLog)

	eventHandler := handlers.NewEventHandler(notificationService, logrusLog)

	//TODO: add exchange to env file
	exchange := viper.GetString("EXCHANGE")
	rabbitMQConsumer := consumer.NewConsumer(conn, exchange, eventHandler, logrusLog)

	if err := rabbitMQConsumer.Start(context.Background()); err != nil {
		logrusLog.Fatalf("Failed to start consumer: %s", err.Error())
	}

	//TODO: replace with proper shutdown handling
	select {}

}
