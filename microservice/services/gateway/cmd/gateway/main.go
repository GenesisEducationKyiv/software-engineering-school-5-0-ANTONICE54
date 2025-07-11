package main

import (
	"weather-forecast/gateway/internal/clients"
	"weather-forecast/gateway/internal/server"
	"weather-forecast/gateway/internal/server/handlers"
	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/proto/subscription"
	"weather-forecast/pkg/proto/weather"

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
		logrusLog.Fatalf("Failed to connect to Weather Service: %v", err)
	}

	defer weatherConn.Close()
	weatherGRPCClient := weather.NewWeatherServiceClient(weatherConn)
	weatherClient := clients.NewWeatherGRPCClient(weatherGRPCClient, logrusLog)
	weatherHandler := handlers.NewWeatherHandler(weatherClient, logrusLog)

	subscConn, err := grpc.NewClient("subscription-service:8082",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logrusLog.Fatalf("Failed to connect to Subscription Service: %v", err)
	}
	defer subscConn.Close()

	subscriptionGRPCClient := subscription.NewSubscriptionServiceClient(subscConn)
	sunbscriptionClient := clients.NewSubscriptionGRPCClient(subscriptionGRPCClient, logrusLog)
	subsbscriptionHandler := handlers.NewSubscriptionHandler(sunbscriptionClient, logrusLog)

	serverPort := viper.GetString("SERVER_PORT")
	app := server.New(weatherHandler, subsbscriptionHandler, logrusLog)
	app.Run(serverPort)
}
