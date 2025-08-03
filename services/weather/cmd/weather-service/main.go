package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"weather-forecast/pkg/logger"
	"weather-service/internal/config"
	"weather-service/internal/domain/usecases"
	"weather-service/internal/infrastructure/cache"
	"weather-service/internal/infrastructure/clients/openweather"
	"weather-service/internal/infrastructure/clients/weatherapi"
	filelogger "weather-service/internal/infrastructure/logger"
	"weather-service/internal/infrastructure/metrics"

	"weather-service/internal/infrastructure/providers"
	"weather-service/internal/infrastructure/providers/roundtrip"
	"weather-service/internal/presentation/server"
	"weather-service/internal/presentation/server/handlers"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to read from config: %s", err.Error())
	}
	logSampler := logger.NewRateSampler(cfg.LogSamplingRate)
	logrusLog := logger.NewLogrus(cfg.ServiceName, cfg.LogLevel, logSampler)

	fileLog, err := filelogger.NewFile(cfg.LogFilePath)
	if err != nil {
		logrusLog.Fatalf("Failed to create file logger: %s", err.Error())
	}
	defer func() {
		if err := fileLog.Close(); err != nil {
			logrusLog.Fatalf("Failed to close log file:%s", err.Error())
		}
	}()

	prometheusMetrics := metrics.NewPrometheus(logrusLog)

	redisCache, err := cache.NewRedis(cfg.RedisSource, logrusLog)
	if err != nil {
		logrusLog.Fatalf("Connect to redis: %s", err.Error())
	}

	providerRoundTrip := roundtrip.New(logrusLog, logrusLog)

	client := http.Client{
		Timeout:   time.Second * 5,
		Transport: providerRoundTrip,
	}

	weatherAPIClient := weatherapi.NewClient(cfg, &client, logrusLog)
	weatherAPIProvider := providers.NewWeatherAPIProvider(weatherAPIClient, logrusLog)
	cacheableWeatherAPIProvider := providers.NewCacheDecorator(weatherAPIProvider, redisCache, prometheusMetrics, logrusLog)

	openWeatherClient := openweather.NewClient(cfg, &client, logrusLog)
	openWeatherProvider := providers.NewOpenWeatherProvider(openWeatherClient, logrusLog)
	cacheableOpenWeatherProvider := providers.NewCacheDecorator(openWeatherProvider, redisCache, prometheusMetrics, logrusLog)

	cacheWeatherProvider := providers.NewCacheWeather(redisCache, prometheusMetrics, logrusLog)

	weatherAPIChainSection := providers.NewWeatherLink(cacheableWeatherAPIProvider)
	openWeatherChainSection := providers.NewWeatherLink(cacheableOpenWeatherProvider)

	cacheWeatherProviderChainSection := providers.NewWeatherLink(cacheWeatherProvider)
	cacheWeatherProviderChainSection.SetNext(weatherAPIChainSection)
	weatherAPIChainSection.SetNext(openWeatherChainSection)

	weatherService := usecases.NewWeatherService(cacheWeatherProviderChainSection, logrusLog)
	weatherHandler := handlers.NewWeatherHandler(weatherService, logrusLog)

	go prometheusMetrics.StartMetricsServer(cfg.MetricsServerPort)

	app := server.New(weatherHandler, logrusLog)
	go func() {
		if err := app.Start(cfg.GRPCPort); err != nil {
			logrusLog.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrusLog.Infof("Shutting down weather service...")
	app.Shutdown()
	logrusLog.Infof("Service stopped gracefully")
}
