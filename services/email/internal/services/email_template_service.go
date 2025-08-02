package services

import (
	"email-service/internal/dto"
	"fmt"
	"weather-forecast/pkg/logger"
)

type (
	Email struct {
		Subject string
		Body    string
	}

	SimpleEmailBuildService struct {
		serverHost string
		logger     logger.Logger
	}
)

func NewSimpleEmailBuild(host string, logger logger.Logger) *SimpleEmailBuildService {
	return &SimpleEmailBuildService{
		serverHost: host,
		logger:     logger,
	}
}

func (s *SimpleEmailBuildService) CreateConfirmationEmail(info *dto.SubscriptionEmailInfo) Email {
	return Email{
		Subject: "Confirm your subscription",
		Body: fmt.Sprintf(
			"You have signed up for an %s newsletter. \nPlease, use this token to confirm your subscription: %s\nOr use this link: %s/confirm/%s",
			info.Frequency, info.Token, s.serverHost, info.Token,
		),
	}
}

func (s *SimpleEmailBuildService) CreateConfirmedEmail(info *dto.ConfirmedEmailInfo) Email {
	return Email{
		Subject: "Subscription confirmed",
		Body: fmt.Sprintf(
			"Congratulations, you have successfully confirmed your %s subscription.\nYou can cancel your subscription using this token: %s\nOr use this link: %s/unsubscribe/%s",
			info.Frequency, info.Token, s.serverHost, info.Token,
		),
	}
}

func (s *SimpleEmailBuildService) CreateUnsubscribeEmail(info *dto.UnsubscribedEmailInfo) Email {
	return Email{
		Subject: "Subscription canceled",
		Body: fmt.Sprintf(
			"You have successfully canceled your %s subscription for city %s.",
			info.Frequency, info.City,
		),
	}
}

func (s *SimpleEmailBuildService) CreateWeatherEmail(info *dto.WeatherSuccess) Email {
	return Email{
		Subject: "Weather Update",
		Body: fmt.Sprintf(
			"Here's the latest weather update for your city: %s\nTemperature: %.1f°C\nHumidity: %d%%\nDescription: %s",
			info.City,
			info.Weather.Temperature,
			info.Weather.Humidity,
			info.Weather.Description,
		),
	}
}

func (s *SimpleEmailBuildService) CreateWeatherErrorEmail(info *dto.WeatherError) Email {
	return Email{
		Subject: "Weather Update",
		Body:    fmt.Sprintf("Sorry, there was an error retrieving weather in your city: %s", info.City),
	}
}
