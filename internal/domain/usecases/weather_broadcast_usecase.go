package usecases

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
	SubscriptionProvider interface {
		ListConfirmedByFrequency(ctx context.Context, frequency models.Frequency, lastID, pageSize int) ([]models.Subscription, error)
	}

	WeatherMailer interface {
		SendWeather(ctx context.Context, email, city string, weather *models.Weather)
		SendError(ctx context.Context, email, city string)
	}

	WeatherBroadcastUseCase struct {
		subscriptionProvider SubscriptionProvider
		weatherService       WeatherProvider
		weatherMailer        WeatherMailer
		logger               logger.Logger
	}
)

func NewWeatherBroadcastUseCase(subscriptionProvider SubscriptionProvider, weatherService WeatherProvider, weatherMailer WeatherMailer, logger logger.Logger) *WeatherBroadcastUseCase {
	return &WeatherBroadcastUseCase{
		subscriptionProvider: subscriptionProvider,
		weatherService:       weatherService,
		weatherMailer:        weatherMailer,
		logger:               logger,
	}
}

func (uc *WeatherBroadcastUseCase) Broadcast(ctx context.Context, frequency models.Frequency) {
	cityWeatherMap := make(map[string]*models.Weather)

	sem := make(chan struct{}, WORKER_AMOUNT)
	wg := &sync.WaitGroup{}

	lastID := 0
	for {
		subscriptions, err := uc.subscriptionProvider.ListConfirmedByFrequency(ctx, frequency, lastID, PAGE_SIZE)
		if err != nil {
			uc.logger.Warnf("Failed to fetch subscriptions: %v", err)
			break
		}
		if len(subscriptions) == 0 {
			break
		}

		for _, subscription := range subscriptions {
			lastID = subscription.ID

			if _, ok := cityWeatherMap[subscription.City]; !ok {
				weather, err := uc.weatherService.GetWeatherByCity(ctx, subscription.City)
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
					uc.weatherMailer.SendWeather(ctx, sub.Email, sub.City, weather)
				} else {
					uc.weatherMailer.SendError(ctx, sub.Email, sub.City)
				}
			}(subscription, cityWeatherMap[subscription.City])
		}
	}
	wg.Wait()
}
