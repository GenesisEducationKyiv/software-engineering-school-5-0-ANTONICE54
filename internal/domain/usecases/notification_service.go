package usecases

import (
	"context"
	"fmt"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/logger"
)

const (
	WEATHER_SUBJECT     = "Weather Update"
	ERROR_BODY_TEMPLATE = "Sorry, there was an error retrieving weather in your city: %s"
)

type (
	Mailer interface {
		Send(ctx context.Context, subject string, body, email string)
	}

	NotificationService struct {
		mailer     Mailer
		serverHost string
		logger     logger.Logger
	}
)

func NewNotificationService(mailer Mailer, serverHost string, logger logger.Logger) *NotificationService {
	return &NotificationService{
		mailer:     mailer,
		serverHost: serverHost,
		logger:     logger,
	}
}

func (s *NotificationService) SendConfirmation(ctx context.Context, email, token string, frequency models.Frequency) {
	subject := "Confirm your subscription"
	body := fmt.Sprintf("You have signed up for an %s newsletter. \n Please, use this token to confirm your subscription: %s\nOr use this link: %s/confirm/%s", frequency, token, s.serverHost, token)
	s.mailer.Send(ctx, subject, body, email)
}

func (s *NotificationService) SendConfirmed(ctx context.Context, email, token string, frequency models.Frequency) {
	subject := "Subscription confirmed"
	body := fmt.Sprintf("Congratulations, you have successfully confirmed your %s subscription.\n You can cancel your subscription using this token: %s\nOr use this link: %s/unsubscribe/%s", frequency, token, s.serverHost, token)
	s.mailer.Send(ctx, subject, body, email)
}

func (s *NotificationService) SendWeather(ctx context.Context, email, city string, weather *models.Weather) {
	body := fmt.Sprintf("Here`s the latest weather update for your city: %s\n Temperature:%.1f C\n Humidity: %d%%\n Description: %s",
		city,
		weather.Temperature,
		weather.Humidity,
		weather.Description,
	)
	s.mailer.Send(ctx, WEATHER_SUBJECT, body, email)
}

func (s *NotificationService) SendError(ctx context.Context, email, city string) {
	body := fmt.Sprintf(ERROR_BODY_TEMPLATE, city)
	s.mailer.Send(ctx, WEATHER_SUBJECT, body, email)
}
