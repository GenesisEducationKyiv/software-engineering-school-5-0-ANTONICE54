package sender

import (
	"context"
	"weather-broadcast-service/internal/dto"

	"weather-broadcast-service/internal/events"
	"weather-forecast/pkg/logger"
)

type (
	EventPublisher interface {
		Publish(ctx context.Context, routingKey string, body []byte) error
	}

	EventSender struct {
		publisher EventPublisher
		logger    logger.Logger
	}
)

func NewEventSender(publisher EventPublisher, logger logger.Logger) *EventSender {
	return &EventSender{
		publisher: publisher,
		logger:    logger,
	}
}

func (s *EventSender) SendWeather(ctx context.Context, info *dto.WeatherMailSuccessInfo) {
	log := s.logger.WithContext(ctx)

	log.Debugf("Creating weather success event: email=%s, city=%s", info.Email, info.City)
	event, err := events.NewWeatherSuccess(info)
	if err != nil {
		log.Errorf("Failed to create weather success event for email %s: %v", info.Email, err)
		return
	}

	routingKey, err := event.RoutingKey()
	if err != nil {
		log.Errorf("Failed to get weather success event routing key for email %s: %v", info.Email, err)
		return
	}

	log.Infof("Publishing weather event: email=%s, city=%s", info.Email, info.City)
	err = s.publisher.Publish(ctx, routingKey, event.Body)
	if err != nil {
		log.Errorf("Failed to publish weather event for email %s: %v", info.Email, err)
		return
	}

	log.Debugf("Weather event published successfully: email=%s", info.Email)

}

func (s *EventSender) SendError(ctx context.Context, info *dto.WeatherMailErrorInfo) {
	log := s.logger.WithContext(ctx)

	log.Debugf("Creating weather error event: email=%s, city=%s", info.Email, info.City)
	event, err := events.NewWeatherError(info)
	if err != nil {
		log.Errorf("Failed to create weather error event for email %s: %v", info.Email, err)
		return
	}

	routingKey, err := event.RoutingKey()
	if err != nil {
		log.Errorf("Failed to get weather error event routing key for email %s: %v", info.Email, err)
		return
	}

	log.Infof("Publishing error weather event: email=%s, city=%s", info.Email, info.City)
	err = s.publisher.Publish(ctx, routingKey, event.Body)
	if err != nil {
		log.Errorf("Failed to publish error weather event for email %s: %v", info.Email, err)
		return
	}

	log.Debugf("Error Weather event published successfully: email=%s", info.Email)

}
