package sender

import (
	"context"
	"weather-broadcast-service/internal/dto"
	"weather-broadcast-service/internal/mappers"

	"weather-forecast/pkg/events"
	"weather-forecast/pkg/logger"
)

type (
	EventPublisher interface {
		Publish(ctx context.Context, event events.Event) error
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

	log.Infof("Publishing weather event: email=%s, city=%s", info.Email, info.City)

	e := mappers.MapWeatherSuccessMailToEvent(info)
	err := s.publisher.Publish(ctx, e)
	if err != nil {
		log.Errorf("Failed to publish weather event for email %s: %v", info.Email, err)

	}

	log.Debugf("Weather event published successfully: email=%s", info.Email)

}
func (s *EventSender) SendError(ctx context.Context, info *dto.WeatherMailErrorInfo) {
	log := s.logger.WithContext(ctx)

	log.Infof("Publishing error weather event: email=%s, city=%s", info.Email, info.City)

	e := mappers.MapWeatherErrorMailToEvent(info)
	err := s.publisher.Publish(ctx, e)
	if err != nil {
		log.Errorf("Failed to publish error weather event for email %s: %v", info.Email, err)
	}
	log.Debugf("Error Weather event published successfully: email=%s", info.Email)

}
