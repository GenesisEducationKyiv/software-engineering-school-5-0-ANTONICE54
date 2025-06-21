package main

import (
	"net/http"
	"time"
	"weather-forecast/internal/domain/usecases"
	"weather-forecast/internal/infrastructure/database"
	"weather-forecast/internal/infrastructure/logger"
	"weather-forecast/internal/infrastructure/mailer"
	"weather-forecast/internal/infrastructure/providers"
	"weather-forecast/internal/infrastructure/repositories"
	"weather-forecast/internal/infrastructure/scheduler"
	"weather-forecast/internal/infrastructure/services"
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

	client := http.Client{
		Timeout: time.Second * 5,
	}

	fileLog, err := logger.NewFile("./logs/weather.log")
	if err != nil {
		logrusLog.Fatalf("Failed to create file logger: %s", err.Error())
	}
	defer func() {
		if err := fileLog.Close(); err != nil {
			logrusLog.Fatalf("Failed to close log file:%s", err.Error())
		}
	}()

	weatherAPIName := viper.GetString("WEATHER_API_NAME")
	weatherAPIURL := viper.GetString("WEATHER_API_URL")
	weatherAPIKey := viper.GetString("WEATHER_API_KEY")
	weatherAPIProvider := providers.NewWeatherAPIProvider(weatherAPIName, weatherAPIURL, weatherAPIKey, &client, logrusLog)
	loggingWeatherAPIProvider := providers.NewLogging(weatherAPIProvider, fileLog)

	openWeatherName := viper.GetString("OPEN_WEATHER_NAME")
	openWeatherURL := viper.GetString("OPEN_WEATHER_URL")
	openWeatherKey := viper.GetString("OPEN_WEATHER_KEY")
	openWeatherProvider := providers.NewOpenWeatherProvider(openWeatherName, openWeatherURL, openWeatherKey, &client, logrusLog)
	loggingOpenWeatherProvider := providers.NewLogging(openWeatherProvider, fileLog)

	weatherAPIChainSection := providers.NewWeatherChain(loggingWeatherAPIProvider)
	openWeatherChainSection := providers.NewWeatherChain(loggingOpenWeatherProvider)
	weatherAPIChainSection.SetNext(openWeatherChainSection)

	weatherService := services.NewWeatherService(weatherAPIChainSection, logrusLog)
	weatherHandler := handlers.NewWeatherHandler(weatherService, logrusLog)

	subscRepo := repositories.NewSubscriptionRepository(db, logrusLog)
	subscUseCase := usecases.NewSubscriptionUseCase(subscRepo, logrusLog)

	serverHost := viper.GetString("SERVER_HOST")
	mailerFrom := viper.GetString("MAILER_FROM")
	mailerHost := viper.GetString("MAILER_HOST")
	mailerPort := viper.GetString("MAILER_PORT")
	mailerUsername := viper.GetString("MAILER_USERNAME")
	mailerPassword := viper.GetString("MAILER_PASSWORD")
	mailer := mailer.NewSMTPMailer(mailerFrom, mailerHost, mailerPort, mailerUsername, mailerPassword, logrusLog)
	notificationService := services.NewNotificationService(mailer, subscUseCase, weatherService, serverHost, logrusLog)

	tokenManager := token.NewUUIDManager()

	subscService := services.NewSubscriptionService(subscUseCase, tokenManager, notificationService, logrusLog)
	subscHandler := handlers.NewSubscriptionHandler(subscService, logrusLog)

	timezone := viper.GetString("TIMEZONE")
	location, err := time.LoadLocation(timezone)
	if err != nil {
		logrusLog.Fatalf("Failed to load timezone: %s", err.Error())
	}

	scheduler := scheduler.New(notificationService, location, logrusLog)

	serverPort := viper.GetString("SERVER_PORT")
	s := server.New(subscHandler, weatherHandler, scheduler, logrusLog)
	s.Run(serverPort)

}
