package main

import (
	"log"
	"subscription-service/internal/domain/usecases"
	"subscription-service/internal/infrastructure/database"
	"subscription-service/internal/infrastructure/publisher"
	"subscription-service/internal/infrastructure/repositories"
	"subscription-service/internal/infrastructure/sender"
	"subscription-service/internal/infrastructure/token"
	"subscription-service/internal/presentation/server"
	"subscription-service/internal/presentation/server/handlers"
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

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}
	defer ch.Close()

	dbHost := viper.GetString("DB_HOST")
	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("DB_NAME")
	dbPort := viper.GetString("DB_PORT")

	db := database.Connect(dbHost, dbUser, dbPassword, dbName, dbPort)
	database.RunMigration(db)

	subscRepo := repositories.NewSubscriptionRepository(db, logrusLog)
	tokenManager := token.NewUUIDManager()

	exchange := viper.GetString("EXCHANGE")
	rabbitMQPublisher := publisher.NewRabbitMQPublisher(ch, exchange, logrusLog)

	eventSender := sender.NewEventSender(rabbitMQPublisher, logrusLog)

	subscUseCase := usecases.NewSubscriptionService(subscRepo, tokenManager, eventSender, logrusLog)
	subscHandler := handlers.NewSubscriptionHandler(subscUseCase, logrusLog)

	app := server.New(subscHandler, logrusLog)
	gRPCPort := viper.GetString("GRPC_PORT")

	app.Start(gRPCPort)
}
