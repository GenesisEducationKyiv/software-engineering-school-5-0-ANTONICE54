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
)

func main() {
	logrusLog := logger.NewLogrus()

	dbHost := "localhost"
	dbUser := "test"
	dbPassword := "test"
	dbName := "test"
	dbPort := "5432"

	db := database.Connect(dbHost, dbUser, dbPassword, dbName, dbPort)
	database.RunMigration(db)

	client := http.Client{
		Timeout: time.Second * 5,
	}
	weatherApiURL := "testAPIURL"
	weatherApiKey := "testAPIKey"
	weatherProvider := providers.NewWeatherProvider(weatherApiURL, weatherApiKey, &client, logrusLog)
	weatherService := services.NewWeatherService(weatherProvider, logrusLog)
	weatherHandler := handlers.NewWeatherHandler(weatherService, logrusLog)

	subscRepo := repositories.NewSubscriptionRepository(db, logrusLog)
	subscUseCase := usecases.NewSubscriptionUseCase(subscRepo, logrusLog)

	serverHost := "http://localhost:8080"
	mailerFrom := "test@test.com"
	mailerHost := "test"
	mailerPort := "test"
	mailerUsername := "test"
	mailerPassword := "test"
	mailer := mailer.NewSMTPMailer(mailerFrom, mailerHost, mailerPort, mailerUsername, mailerPassword, logrusLog)
	notificationService := services.NewNotificationService(mailer, subscUseCase, weatherService, serverHost, logrusLog)

	tokenManager := token.NewUUIDManager()

	subscService := services.NewSubscriptionService(subscUseCase, tokenManager, notificationService, logrusLog)
	subscHandler := handlers.NewSubscriptionHandler(subscService, logrusLog)

	timezone := "Europe/Kyiv"
	location, err := time.LoadLocation(timezone)
	if err != nil {
		logrusLog.Fatalf("Failed to load timezone: %s", err.Error())
	}

	scheduler := scheduler.New(notificationService, location, logrusLog)

	serverPort := "8080"
	s := server.New(subscHandler, weatherHandler, scheduler, logrusLog)
	s.Run(serverPort)
}
