package services

import (
	"context"
	"email-service/internal/dto"
	"weather-forecast/pkg/logger"
)

type (
	EmailBuildService interface {
		CreateConfirmationEmail(info *dto.SubscriptionEmailInfo) Email
		CreateConfirmedEmail(info *dto.ConfirmedEmailInfo) Email
		CreateWeatherEmail(info *dto.WeatherSuccess) Email
		CreateWeatherErrorEmail(info *dto.WeatherError) Email
	}

	Mailer interface {
		Send(ctx context.Context, subject string, body, email string)
	}

	NotificationService struct {
		mailer       Mailer
		emailBuilder EmailBuildService
		logger       logger.Logger
	}
)

func NewNotificationService(mailer Mailer, emailBuilder EmailBuildService, logger logger.Logger) *NotificationService {
	return &NotificationService{
		mailer:       mailer,
		emailBuilder: emailBuilder,
		logger:       logger,
	}
}

func (s *NotificationService) SendConfirmation(ctx context.Context, info *dto.SubscriptionEmailInfo) {
	email := s.emailBuilder.CreateConfirmationEmail(info)
	s.mailer.Send(ctx, email.Subject, email.Body, info.Email)
}

func (s *NotificationService) SendConfirmed(ctx context.Context, info *dto.ConfirmedEmailInfo) {
	email := s.emailBuilder.CreateConfirmedEmail(info)
	s.mailer.Send(ctx, email.Subject, email.Body, info.Email)
}

func (s *NotificationService) SendWeather(ctx context.Context, info *dto.WeatherSuccess) {
	email := s.emailBuilder.CreateWeatherEmail(info)
	s.mailer.Send(ctx, email.Subject, email.Body, info.Email)
}

func (s *NotificationService) SendError(ctx context.Context, info *dto.WeatherError) {
	email := s.emailBuilder.CreateWeatherErrorEmail(info)
	s.mailer.Send(ctx, email.Subject, email.Body, info.Email)
}
