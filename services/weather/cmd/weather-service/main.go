package main

import (
	"net/http"
	"time"
	"weather-forecast/pkg/logger"
	"weather-service/internal/config"
	"weather-service/internal/domain/usecases"
	"weather-service/internal/infrastructure/cache"
	filelogger "weather-service/internal/infrastructure/logger"
	"weather-service/internal/infrastructure/metrics"

	"weather-service/internal/infrastructure/providers"
	"weather-service/internal/infrastructure/providers/roundtrip"
	"weather-service/internal/presentation/server"
	"weather-service/internal/presentation/server/handlers"
)

func main() {

	logrusLog := logger.NewLogrus()

	cfg, err := config.Load(logrusLog)
	if err != nil {
		logrusLog.Fatalf("Failed to read from config: %s", err.Error())
	}

	fileLog, err := filelogger.NewFile(cfg.LogFilePath)
	if err != nil {
		logrusLog.Fatalf("Failed to create file logger: %s", err.Error())
	}
	defer func() {
		if err := fileLog.Close(); err != nil {
			logrusLog.Fatalf("Failed to close log file:%s", err.Error())
		}
	}()

	prometherusMetrics := metrics.NewPrometheus(logrusLog)

	redisCache, err := cache.NewRedis(cfg.RedisSource, logrusLog)
	if err != nil {
		logrusLog.Fatalf("Connect to redis: %s", err.Error())
	}

	providerRoundTrip := roundtrip.New(fileLog, logrusLog)

	client := http.Client{
		Timeout:   time.Second * 5,
		Transport: providerRoundTrip,
	}

	weatherAPIProvider := providers.NewWeatherAPIProvider(cfg, &client, logrusLog)
	cacheableWeatherAPIProvider := providers.NewCacheDecorator(weatherAPIProvider, redisCache, prometherusMetrics, logrusLog)

	openWeatherProvider := providers.NewOpenWeatherProvider(cfg, &client, logrusLog)
	cacheableOpenWeatherProvider := providers.NewCacheDecorator(openWeatherProvider, redisCache, prometherusMetrics, logrusLog)

	cacheWeatherProvider := providers.NewCacheWeather(redisCache, prometherusMetrics, logrusLog)

	weatherAPIChainSection := providers.NewWeatherLink(cacheableWeatherAPIProvider)
	openWeatherChainSection := providers.NewWeatherLink(cacheableOpenWeatherProvider)

	cacheWeatherProviderChainSection := providers.NewWeatherLink(cacheWeatherProvider)
	cacheWeatherProviderChainSection.SetNext(weatherAPIChainSection)
	weatherAPIChainSection.SetNext(openWeatherChainSection)

	weatherService := usecases.NewWeatherService(cacheWeatherProviderChainSection, logrusLog)
	weatherHandler := handlers.NewWeatherHandler(weatherService, logrusLog)

	app := server.New(weatherHandler, logrusLog)

	go prometherusMetrics.StartMetricsServer(cfg.MetricsServerPort)

	if err := app.Start(cfg.GRPCPort); err != nil {
		logrusLog.Fatalf("Failed to start server: %v", err)
	}
}
