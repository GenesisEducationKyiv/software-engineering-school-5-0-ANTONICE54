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
	e := mappers.MapWeatherSuccessMailToEvent(info)
	err := s.publisher.Publish(ctx, e)
	if err != nil {
		s.logger.Warnf("failed to publish event: %s", err.Error())
	}
}
func (s *EventSender) SendError(ctx context.Context, info *dto.WeatherMailErrorInfo) {
	e := mappers.MapWeatherErrorMailToEvent(info)
	err := s.publisher.Publish(ctx, e)
	if err != nil {
		s.logger.Warnf("failed to publish event: %s", err.Error())
	}
}
