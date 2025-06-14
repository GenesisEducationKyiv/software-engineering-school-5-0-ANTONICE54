package services

import (
	"context"
	"sync"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/logger"
)

const (
	PAGE_SIZE     = 100
	WORKER_AMOUNT = 10
)

type (
	WeatherServiceI interface {
		GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error)
	}

	ListSubscriptionUseCase interface {
		ListByFrequency(ctx context.Context, frequency models.Frequency, lastID, pageSize int) ([]models.Subscription, error)
	}

	WeatherMailer interface {
		SendWeather(ctx context.Context, email, city string, weather *models.Weather)
		SendError(ctx context.Context, email, city string)
	}

	WeatherBroadcastService struct {
		subscriptionUseCase ListSubscriptionUseCase
		weatherService      WeatherServiceI
		weatherMailer       WeatherMailer
		logger              logger.Logger
	}
)

func NewWeatherBroadcastService(subscriptionUC ListSubscriptionUseCase, weatherService WeatherServiceI, weatherMailer WeatherMailer, logger logger.Logger) *WeatherBroadcastService {
	return &WeatherBroadcastService{
		subscriptionUseCase: subscriptionUC,
		weatherService:      weatherService,
		weatherMailer:       weatherMailer,
		logger:              logger,
	}
}

func (s *WeatherBroadcastService) Broadcast(ctx context.Context, frequency models.Frequency) {
	cityWeatherMap := make(map[string]*models.Weather)

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

			if _, ok := cityWeatherMap[subscription.City]; !ok {
				weather, err := s.weatherService.GetWeatherByCity(ctx, subscription.City)
				if err != nil {
					cityWeatherMap[subscription.City] = nil
				}
				cityWeatherMap[subscription.City] = weather
			}

			sem <- struct{}{}
			wg.Add(1)

			go func(sub models.Subscription, weather *models.Weather) {
				defer func() { <-sem }()
				defer wg.Done()
				if weather != nil {
					s.weatherMailer.SendWeather(ctx, sub.Email, sub.City, weather)
				} else {
					s.weatherMailer.SendError(ctx, sub.Email, sub.City)
				}
			}(subscription, cityWeatherMap[subscription.City])
		}
	}
	wg.Wait()
}
