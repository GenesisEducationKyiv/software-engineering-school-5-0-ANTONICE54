package services

import (
	"context"
	"sync"
	"weather-broadcast-service/internal/dto"
	"weather-broadcast-service/internal/models"
	"weather-forecast/pkg/logger"
)

const (
	PAGE_SIZE     = 100
	WORKER_AMOUNT = 10
)

type (
	WeatherClient interface {
		GetWeatherByCity(ctx context.Context, city string) (*dto.Weather, error)
	}

	SubscriptionClient interface {
		ListByFrequency(ctx context.Context, query dto.ListSubscriptionsQuery) (*dto.SubscriptionList, error)
	}

	WeatherMailer interface {
		SendWeather(ctx context.Context, info *dto.WeatherMailSuccessInfo)
		SendError(ctx context.Context, info *dto.WeatherMailErrorInfo)
	}

	WeatherBroadcastService struct {
		subscriptionClient SubscriptionClient
		weatherClient      WeatherClient
		weatherMailer      WeatherMailer
		logger             logger.Logger
	}
)

func NewWeatherBroadcastService(subscriptionClient SubscriptionClient, weatherClient WeatherClient, weatherMailer WeatherMailer, logger logger.Logger) *WeatherBroadcastService {
	return &WeatherBroadcastService{
		subscriptionClient: subscriptionClient,
		weatherClient:      weatherClient,
		weatherMailer:      weatherMailer,
		logger:             logger,
	}
}

func (s *WeatherBroadcastService) Broadcast(ctx context.Context, frequency models.Frequency) {
	log := s.logger.WithContext(ctx)

	log.Debugf("Starting broadcast process for %s subscription", frequency)

	cityWeatherMap := make(map[string]*dto.Weather)
	sem := make(chan struct{}, WORKER_AMOUNT)
	wg := &sync.WaitGroup{}

	lastID := 0
	for {

		query := dto.ListSubscriptionsQuery{
			Frequency: frequency,
			LastID:    lastID,
			PageSize:  PAGE_SIZE,
		}

		log.Debugf("Getting subsctiption list batch from index %d with page size=%d", query.LastID, query.PageSize)
		res, err := s.subscriptionClient.ListByFrequency(ctx, query)
		if err != nil {
			log.Errorf("Failed to fetch subscriptions for %s broadcast: %v", frequency, err)
			break
		}

		subscriptions := res.Subscriptions
		lastID = res.LastIndex

		if len(subscriptions) == 0 {
			log.Debugf("No more subscriptions found, stopping broadcast")
			break
		}

		for _, subscription := range subscriptions {

			if _, ok := cityWeatherMap[subscription.City]; !ok {
				log.Debugf("Getting weather for new city: %s", subscription.City)

				weather, err := s.weatherClient.GetWeatherByCity(ctx, subscription.City)
				log.Infof("Weather: %v", weather)

				if err != nil {
					log.Warnf("Failed to get weather for city %s: %v", subscription.City, err)
					cityWeatherMap[subscription.City] = nil
				} else {
					log.Debugf("Weather fetched successfully for city: %s", subscription.City)
					cityWeatherMap[subscription.City] = weather
				}

			}

			sem <- struct{}{}
			wg.Add(1)

			go func(sub dto.Subscription, weather *dto.Weather) {
				defer func() { <-sem }()
				defer wg.Done()
				if weather != nil {
					log.Debugf("Sending weather email to: %s for city: %s", sub.Email, sub.City)
					info := &dto.WeatherMailSuccessInfo{
						Email:   sub.Email,
						City:    sub.City,
						Weather: *weather,
					}

					s.weatherMailer.SendWeather(ctx, info)
				} else {
					log.Debugf("Sending error email to: %s for city: %s", sub.Email, sub.City)
					info := &dto.WeatherMailErrorInfo{
						Email: sub.Email,
						City:  sub.City,
					}

					s.weatherMailer.SendError(ctx, info)
				}
			}(subscription, cityWeatherMap[subscription.City])
		}
	}
	wg.Wait()
}
