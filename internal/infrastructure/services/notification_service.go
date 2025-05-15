package services

import (
	"fmt"
	"sync"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/logger"
)

type (
	MailerI interface {
		Send(subject string, body, email string)
	}
	ListSubscriptionUseCaseI interface {
		ListByFrequency(frequency models.Frequency) ([]models.Subscription, error)
	}
	WeatherServiceI interface {
		GetWeatherByCity(city string) (*models.Weather, error)
	}
	NotificationService struct {
		mailer              MailerI
		subscriptionUseCase ListSubscriptionUseCaseI
		weatherService      WeatherServiceI
		serverHost          string
		logger              logger.Logger
	}
)

func NewNotificationService(mailer MailerI, subscriptionUC ListSubscriptionUseCaseI, weatherService WeatherServiceI, serverHost string, logger logger.Logger) *NotificationService {
	return &NotificationService{
		mailer:              mailer,
		subscriptionUseCase: subscriptionUC,
		weatherService:      weatherService,
		serverHost:          serverHost,
		logger:              logger,
	}
}

func (s *NotificationService) SendConfirmation(email, token string, frequency models.Frequency) {
	subject := "Confirm your subscription"
	body := fmt.Sprintf("You have signed up for an %s newsletter. \n Please, use this token to confirm your subscription: %s\nOr use this link: %s/confirm/%s", frequency, token, s.serverHost, token)
	s.mailer.Send(subject, body, email)
}

func (s *NotificationService) SendConfirmed(email, token string, frequency models.Frequency) {
	subject := "Subscription confirmed"
	body := fmt.Sprintf("Congratulations, you have successfully confirmed your %s subscription.\n You can cancel your subscription using this token: %s\nOr use this link: %s/unsubscribe/%s", frequency, token, s.serverHost, token)
	s.mailer.Send(subject, body, email)
}

func (s *NotificationService) SendWeather(frequency models.Frequency) {
	subscriptions, err := s.subscriptionUseCase.ListByFrequency(frequency)
	if err != nil {
		s.logger.Warnf("Failed to send %s weather update", frequency)
		return
	}

	wg := &sync.WaitGroup{}
	subject := "Weather update"
	for _, subscription := range subscriptions {
		wg.Add(1)
		go func(email, city string) {
			defer wg.Done()
			weather, err := s.weatherService.GetWeatherByCity(city)
			var body string
			if err != nil {
				body = "Sorry, there was an error retrieving weather in your city: " + city
			} else {
				body = fmt.Sprintf("Here`s the latest weather update for your city: %s\n Temperature:%.1f C\n Humidity: %d%%\n Description: %s",
					city,
					weather.Temperature,
					weather.Humidity,
					weather.Description,
				)
			}
			s.mailer.Send(subject, body, email)
		}(subscription.Email, subscription.City)
	}
	wg.Wait()
}
