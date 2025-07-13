package main

import (
	"net/http"
	"time"
	"weather-forecast/pkg/logger"
	"weather-service/internal/domain/usecases"
	"weather-service/internal/infrastructure/cache"
	filelogger "weather-service/internal/infrastructure/logger"
	"weather-service/internal/infrastructure/metrics"

	"weather-service/internal/infrastructure/providers"
	"weather-service/internal/infrastructure/providers/roundtrip"
	"weather-service/internal/presentation/server"
	"weather-service/internal/presentation/server/handlers"

	"github.com/spf13/viper"
)

func main() {

	logrusLog := logger.NewLogrus()

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		logrusLog.Fatalf("Failed to read from config: %s", err.Error())
	}

	fileLog, err := filelogger.NewFile("./logs/weather.txt")
	if err != nil {
		logrusLog.Fatalf("Failed to create file logger: %s", err.Error())
	}
	defer func() {
		if err := fileLog.Close(); err != nil {
			logrusLog.Fatalf("Failed to close log file:%s", err.Error())
		}
	}()

	prometherusMetrics := metrics.NewPrometheus(logrusLog)

	residSource := viper.GetString("REDIS_SOURCE")
	redisCache, err := cache.NewRedis(residSource, logrusLog)
	if err != nil {
		logrusLog.Fatalf("Connect to redis: %s", err.Error())
	}

	providerRoundTrip := roundtrip.New(fileLog, logrusLog)

	client := http.Client{
		Timeout:   time.Second * 5,
		Transport: providerRoundTrip,
	}

	weatherAPIURL := viper.GetString("WEATHER_API_URL")
	weatherAPIKey := viper.GetString("WEATHER_API_KEY")
	weatherAPIProvider := providers.NewWeatherAPIProvider(weatherAPIURL, weatherAPIKey, &client, logrusLog)
	cacheableWeatherAPIProvider := providers.NewCacheDecorator(weatherAPIProvider, redisCache, prometherusMetrics, logrusLog)

	openWeatherURL := viper.GetString("OPEN_WEATHER_URL")
	openWeatherKey := viper.GetString("OPEN_WEATHER_KEY")
	openWeatherProvider := providers.NewOpenWeatherProvider(openWeatherURL, openWeatherKey, &client, logrusLog)
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

	metricsServerPort := viper.GetString("METRICS_SERVER_PORT")
	go prometherusMetrics.StartMetricsServer(metricsServerPort)

	gRPCPort := viper.GetString("GRPC_PORT")

	app.Start(gRPCPort)

}
