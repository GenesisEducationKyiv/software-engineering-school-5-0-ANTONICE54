package main

import (
	"context"
	"net/http"
	"time"
	"weather-forecast/internal/domain/usecases"
	"weather-forecast/internal/infrastructure/cache"
	"weather-forecast/internal/infrastructure/database"
	"weather-forecast/internal/infrastructure/logger"
	"weather-forecast/internal/infrastructure/mailer"
	"weather-forecast/internal/infrastructure/metrics"
	"weather-forecast/internal/infrastructure/providers"
	"weather-forecast/internal/infrastructure/providers/roundtrip"
	"weather-forecast/internal/infrastructure/repositories"
	"weather-forecast/internal/infrastructure/scheduler"
	"weather-forecast/internal/infrastructure/token"
	"weather-forecast/internal/presentation/server"
	"weather-forecast/internal/presentation/server/handlers"

	"github.com/spf13/viper"
)

func main() {

	logrusLog := logger.NewLogrus()

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		logrusLog.Fatalf("Failed to read from config: %s", err.Error())
	}

	dbHost := viper.GetString("DB_HOST")
	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("DB_NAME")
	dbPort := viper.GetString("DB_PORT")

	db := database.Connect(dbHost, dbUser, dbPassword, dbName, dbPort)
	database.RunMigration(db)

	fileLog, err := logger.NewFile("./logs/weather.txt")
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

	subscRepo := repositories.NewSubscriptionRepository(db, logrusLog)

	serverHost := viper.GetString("SERVER_HOST")
	mailerFrom := viper.GetString("MAILER_FROM")
	mailerHost := viper.GetString("MAILER_HOST")
	mailerPort := viper.GetString("MAILER_PORT")
	mailerUsername := viper.GetString("MAILER_USERNAME")
	mailerPassword := viper.GetString("MAILER_PASSWORD")
	mailer := mailer.NewSMTPMailer(mailerFrom, mailerHost, mailerPort, mailerUsername, mailerPassword, logrusLog)
	notificationService := usecases.NewNotificationService(mailer, serverHost, logrusLog)

	tokenManager := token.NewUUIDManager()

	subscService := usecases.NewSubscriptionService(subscRepo, tokenManager, notificationService, logrusLog)
	subscHandler := handlers.NewSubscriptionHandler(subscService, logrusLog)

	timezone := viper.GetString("TIMEZONE")
	location, err := time.LoadLocation(timezone)
	if err != nil {
		logrusLog.Fatalf("Failed to load timezone: %s", err.Error())
	}

	weatherBroadcastService := usecases.NewWeatherBroadcastService(subscRepo, weatherService, notificationService, logrusLog)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler := scheduler.New(ctx, weatherBroadcastService, location, logrusLog)

	serverPort := viper.GetString("SERVER_PORT")
	s := server.New(subscHandler, weatherHandler, scheduler, logrusLog)

	metricsServerPort := viper.GetString("METRICS_SERVER_PORT")
	go prometherusMetrics.StartMetricsServer(metricsServerPort)

	s.Run(serverPort)

}
