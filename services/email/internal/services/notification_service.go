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
		CreateUnsubscribeEmail(info *dto.UnsubscribedEmailInfo) Email
		CreateWeatherEmail(info *dto.WeatherSuccess) Email
		CreateWeatherErrorEmail(info *dto.WeatherError) Email
	}

	Mailer interface {
		Send(ctx context.Context, subject string, body, email string) error
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
	log := s.logger.WithContext(ctx)

	email := s.emailBuilder.CreateConfirmationEmail(info)
	log.Debugf("Created subscription email with subject: '%s' for %s", email.Subject, info.Email)

	err := s.mailer.Send(ctx, email.Subject, email.Body, info.Email)
	if err != nil {
		log.Errorf("Failed to send confirmation email to %s: %v", info.Email, err)
	} else {
		log.Infof("Confirmation email sent successfully to %s", info.Email)
	}
}

func (s *NotificationService) SendConfirmed(ctx context.Context, info *dto.ConfirmedEmailInfo) {
	log := s.logger.WithContext(ctx)

	email := s.emailBuilder.CreateConfirmedEmail(info)
	log.Debugf("Created confirmation email with subject: '%s' for %s", email.Subject, info.Email)

	err := s.mailer.Send(ctx, email.Subject, email.Body, info.Email)
	if err != nil {
		log.Errorf("Failed to send confirmed email to %s: %v", info.Email, err)
	} else {
		log.Infof("Confirmed email sent successfully to %s", info.Email)
	}
}

func (s *NotificationService) SendUnsubscribed(ctx context.Context, info *dto.UnsubscribedEmailInfo) {
	log := s.logger.WithContext(ctx)

	email := s.emailBuilder.CreateUnsubscribeEmail(info)
	log.Debugf("Created unsubscribed email with subject: '%s' for %s", email.Subject, info.Email)

	err := s.mailer.Send(ctx, email.Subject, email.Body, info.Email)
	if err != nil {
		log.Errorf("Failed to send canceled email to %s: %v", info.Email, err)
	} else {
		log.Infof("Canceled email sent successfully to %s", info.Email)
	}
}

func (s *NotificationService) SendWeather(ctx context.Context, info *dto.WeatherSuccess) {
	log := s.logger.WithContext(ctx)

	email := s.emailBuilder.CreateWeatherEmail(info)
	log.Debugf("Created weather email for: %s, city: %s", info.Email, info.City)

	err := s.mailer.Send(ctx, email.Subject, email.Body, info.Email)
	if err != nil {
		log.Errorf("Failed to send weather email to %s (city: %s): %v", info.Email, info.City, err)
	} else {
		log.Infof("Weather email sent successfully to %s (city: %s)", info.Email, info.City)
	}
}

func (s *NotificationService) SendError(ctx context.Context, info *dto.WeatherError) {
	log := s.logger.WithContext(ctx)

	email := s.emailBuilder.CreateWeatherErrorEmail(info)
	log.Debugf("Created weather error email for: %s, city: %s", info.Email, info.City)

	err := s.mailer.Send(ctx, email.Subject, email.Body, info.Email)
	if err != nil {
		log.Errorf("Failed to send weather error email to %s (city: %s): %v", info.Email, info.City, err)
	} else {
		log.Infof("Weather error email sent successfully to %s (city: %s)", info.Email, info.City)
	}
}
