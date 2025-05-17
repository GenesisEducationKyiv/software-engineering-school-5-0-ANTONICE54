package services

import (
	"context"
	"fmt"
	"sync"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/logger"
)

const (
	WEATHER_SUBJECT     = "Weather Update"
	ERROR_BODY_TEMPLATE = "Sorry, there was an error retrieving weather in your city: %s"
	PAGE_SIZE           = 100
	WORKER_AMOUNT       = 10
)

type (
	Mailer interface {
		Send(subject string, body, email string)
	}
	ListSubscriptionUseCase interface {
		ListByFrequency(ctx context.Context, frequency models.Frequency, lastID, pageSize int) ([]models.Subscription, error)
	}
	WeatherServiceI interface {
		GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error)
	}
	NotificationService struct {
		mailer              Mailer
		subscriptionUseCase ListSubscriptionUseCase
		weatherService      WeatherServiceI
		serverHost          string
		logger              logger.Logger
	}
)

func NewNotificationService(mailer Mailer, subscriptionUC ListSubscriptionUseCase, weatherService WeatherServiceI, serverHost string, logger logger.Logger) *NotificationService {
	return &NotificationService{
		mailer:              mailer,
		subscriptionUseCase: subscriptionUC,
		weatherService:      weatherService,
		serverHost:          serverHost,
		logger:              logger,
	}
}

func (s *NotificationService) SendConfirmation(_ context.Context, email, token string, frequency models.Frequency) {
	subject := "Confirm your subscription"
	body := fmt.Sprintf("You have signed up for an %s newsletter. \n Please, use this token to confirm your subscription: %s\nOr use this link: %s/confirm/%s", frequency, token, s.serverHost, token)
	s.mailer.Send(subject, body, email)
}

func (s *NotificationService) SendConfirmed(_ context.Context, email, token string, frequency models.Frequency) {
	subject := "Subscription confirmed"
	body := fmt.Sprintf("Congratulations, you have successfully confirmed your %s subscription.\n You can cancel your subscription using this token: %s\nOr use this link: %s/unsubscribe/%s", frequency, token, s.serverHost, token)
	s.mailer.Send(subject, body, email)
}

func (s *NotificationService) SendWeather(ctx context.Context, frequency models.Frequency) {

	sem := make(chan struct{}, WORKER_AMOUNT)
	wg := &sync.WaitGroup{}
	lastID := 0
	for {
		subscriptions, err := s.subscriptionUseCase.ListByFrequency(ctx, frequency, lastID, PAGE_SIZE)
		if err != nil {
			s.logger.Warnf("Failed to fetch subscriptions: %v", err)
			break
		}

		if len(subscriptions) == 0 {
			break
		}

		for _, subscription := range subscriptions {
			lastID = subscription.ID

			sem <- struct{}{}
			wg.Add(1)

			go func(sub models.Subscription) {
				defer func() { <-sem }()

				s.processWeatherEmail(ctx, sub.Email, sub.City, wg)
			}(subscription)
		}
	}
	wg.Wait()
}

func (s *NotificationService) processWeatherEmail(ctx context.Context, email, city string, wg *sync.WaitGroup) {
	defer wg.Done()
	weather, err := s.weatherService.GetWeatherByCity(ctx, city)
	var body string
	if err != nil {
		body = fmt.Sprintf(ERROR_BODY_TEMPLATE, city)
	} else {
		body = fmt.Sprintf("Here`s the latest weather update for your city: %s\n Temperature:%.1f C\n Humidity: %d%%\n Description: %s",
			city,
			weather.Temperature,
			weather.Humidity,
			weather.Description,
		)
	}
	s.mailer.Send(WEATHER_SUBJECT, body, email)
}
