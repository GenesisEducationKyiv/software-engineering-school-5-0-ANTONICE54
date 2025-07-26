package sender

import (
	"context"
	"weather-broadcast-service/internal/dto"

	"weather-broadcast-service/internal/events"
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

	event, err := events.NewWeatherSuccess(info)
	if err != nil {
		s.logger.Warnf("failed to create event: %s", err.Error())
		return
	}

	err = s.publisher.Publish(ctx, *event)
	if err != nil {
		s.logger.Warnf("failed to publish event: %s", err.Error())
	}
}
func (s *EventSender) SendError(ctx context.Context, info *dto.WeatherMailErrorInfo) {
	event, err := events.NewWeatherError(info)
	if err != nil {
		s.logger.Warnf("failed to create event: %s", err.Error())
		return
	}
	err = s.publisher.Publish(ctx, *event)
	if err != nil {
		s.logger.Warnf("failed to publish event: %s", err.Error())
	}
}
