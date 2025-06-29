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

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func main() {
	logrusLog := logger.NewLogrus()

	viper.SetConfigFile("tests/e2e/testapp/test.env")
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
	weatherApiURL := viper.GetString("WEATHER_API_URL")
	weatherApiKey := viper.GetString("WEATHER_API_KEY")
	weatherProvider := providers.NewWeatherProvider(weatherApiURL, weatherApiKey, &client, logrusLog)
	weatherService := services.NewWeatherService(weatherProvider, logrusLog)
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
	notificationService := services.NewNotificationService(mailer, serverHost, logrusLog)

	tokenManager := token.NewUUIDManager()

	subscService := services.NewSubscriptionService(subscUseCase, tokenManager, notificationService, logrusLog)
	subscHandler := handlers.NewSubscriptionHandler(subscService, logrusLog)

	timezone := viper.GetString("TIMEZONE")
	location, err := time.LoadLocation(timezone)
	if err != nil {
		logrusLog.Fatalf("Failed to load timezone: %s", err.Error())
	}

	weatherBroadcastService := services.NewWeatherBroadcastService(subscUseCase, weatherService, notificationService, logrusLog)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	scheduler := scheduler.New(ctx, weatherBroadcastService, location, logrusLog)

	serverPort := viper.GetString("SERVER_PORT")
	s := server.New(subscHandler, weatherHandler, scheduler, logrusLog)
	go startDebugServer(db, logrusLog)
	s.Run(serverPort)
}

func startDebugServer(db *gorm.DB, logger logger.Logger) {
	debugRouter := gin.New()

	debugRouter.POST("/clear", func(c *gin.Context) {
		err := db.Exec("TRUNCATE TABLE subscriptions").Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "database cleared"})
	})

	port := viper.GetString("DEBUG_SERVER_PORT")
	logger.Infof("Starting debug server on %s...", port)
	if err := debugRouter.Run(":" + port); err != nil {
		logger.Fatalf("Debug server failed: %s", err)
	}
}
