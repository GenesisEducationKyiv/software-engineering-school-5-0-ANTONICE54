package main

import (
	"log"
	"weather-forecast/gateway/internal/clients"
	"weather-forecast/gateway/internal/config"
	"weather-forecast/gateway/internal/metrics"
	"weather-forecast/gateway/internal/server"
	"weather-forecast/gateway/internal/server/handlers"
	grpcpkg "weather-forecast/pkg/grpc"
	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/proto/subscription"
	"weather-forecast/pkg/proto/weather"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to read from config: %s", err.Error())
	}
	logSampler := logger.NewRateSampler(cfg.LogSamplingRate)
	logrusLog := logger.NewLogrus(cfg.ServiceName, cfg.LogLevel, logSampler)
	prometheusMetrics := metrics.NewPrometheus(logrusLog)

	weatherConn, err := grpcpkg.ConnectWithRetry(cfg.WeatherServiceAddress, cfg.GRPC, logrusLog)

	if err != nil {
		logrusLog.Fatalf("Failed to connect to Weather Service: %v", err)
	}
	defer func() {
		if err := weatherConn.Close(); err != nil {
			logrusLog.Errorf("Failed to close gRPC connection with weather service: %v", err)
		}
	}()

	weatherGRPCClient := weather.NewWeatherServiceClient(weatherConn)
	weatherClient := clients.NewWeatherGRPCClient(weatherGRPCClient, logrusLog)
	weatherHandler := handlers.NewWeatherHandler(weatherClient, logrusLog)

	subscConn, err := grpcpkg.ConnectWithRetry(cfg.SubscriptionServiceAddress, cfg.GRPC, logrusLog)

	if err != nil {
		logrusLog.Fatalf("Failed to connect to Subscription Service: %v", err)
	}
	defer func() {
		if err := subscConn.Close(); err != nil {
			logrusLog.Errorf("Failed to close gRPC connection with subscription service: %v", err)
		}
	}()

	subscriptionGRPCClient := subscription.NewSubscriptionServiceClient(subscConn)
	subscriptionClient := clients.NewSubscriptionGRPCClient(subscriptionGRPCClient, logrusLog)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionClient, logrusLog)

	app := server.New(weatherHandler, subscriptionHandler, prometheusMetrics, logrusLog)

	go prometheusMetrics.StartMetricsServer(cfg.MetricsServerPort)

	app.Run(cfg.ServerPort)
}
