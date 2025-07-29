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
	err := s.mailer.Send(ctx, email.Subject, email.Body, info.Email)
	if err != nil {
		log.Debugf("Failed to send email with subject %s to %s. Due to error: %s", email.Subject, info.Email, err.Error())
	} else {
		log.Debugf("Email sent successfully to %s", info.Email)
	}
}

func (s *NotificationService) SendConfirmed(ctx context.Context, info *dto.ConfirmedEmailInfo) {
	log := s.logger.WithContext(ctx)

	email := s.emailBuilder.CreateConfirmedEmail(info)
	err := s.mailer.Send(ctx, email.Subject, email.Body, info.Email)
	if err != nil {
		log.Debugf("Failed to send email with subject %s to %s. Due to error: %s", email.Subject, info.Email, err.Error())
	} else {
		log.Debugf("Email sent successfully to %s", info.Email)
	}
}

func (s *NotificationService) SendUnsubscribed(ctx context.Context, info *dto.UnsubscribedEmailInfo) {
	log := s.logger.WithContext(ctx)

	email := s.emailBuilder.CreateUnsubscribeEmail(info)
	err := s.mailer.Send(ctx, email.Subject, email.Body, info.Email)
	if err != nil {
		log.Debugf("Failed to send email with subject %s to %s. Due to error: %s", email.Subject, info.Email, err.Error())
	} else {
		log.Debugf("Email sent successfully to %s", info.Email)
	}
}

func (s *NotificationService) SendWeather(ctx context.Context, info *dto.WeatherSuccess) {
	log := s.logger.WithContext(ctx)

	email := s.emailBuilder.CreateWeatherEmail(info)
	err := s.mailer.Send(ctx, email.Subject, email.Body, info.Email)
	if err != nil {
		log.Debugf("Failed to send email with subject %s to %s. Due to error: %s", email.Subject, info.Email, err.Error())
	} else {
		log.Debugf("Email sent successfully to %s", info.Email)
	}
}

func (s *NotificationService) SendError(ctx context.Context, info *dto.WeatherError) {
	log := s.logger.WithContext(ctx)

	email := s.emailBuilder.CreateWeatherErrorEmail(info)
	err := s.mailer.Send(ctx, email.Subject, email.Body, info.Email)
	if err != nil {
		log.Debugf("Failed to send email with subject %s to %s. Due to error: %s", email.Subject, info.Email, err.Error())
	} else {
		log.Debugf("Email sent successfully to %s", info.Email)
	}
}
