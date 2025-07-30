package processors

import (
	"context"
	"email-service/internal/dto"
	"email-service/internal/mappers"
	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/proto/events"

	"google.golang.org/protobuf/proto"
)

const (
	ConfirmationRoute   = "emails.subscription"
	ConfirmedRoute      = "emails.confirmed"
	WeatherSuccessRoute = "emails.weather.success"
	WeatherErrorRoute   = "emails.weather.error"
)

type (
	NotificationService interface {
		SendConfirmation(ctx context.Context, info *dto.SubscriptionEmailInfo)
		SendConfirmed(ctx context.Context, info *dto.ConfirmedEmailInfo)
		SendWeather(ctx context.Context, info *dto.WeatherSuccess)
		SendError(ctx context.Context, info *dto.WeatherError)
	}

	EventProcessor struct {
		sender NotificationService
		logger logger.Logger
	}
)

func NewEventProcessor(sender NotificationService, logger logger.Logger) *EventProcessor {
	return &EventProcessor{
		sender: sender,
		logger: logger,
	}
}

func (h *EventProcessor) Handle(ctx context.Context, routingKey string, body []byte) {

	switch routingKey {

	case ConfirmationRoute:
		e := &events.SubscriptionEvent{}
		if err := proto.Unmarshal(body, e); err != nil {
			h.logger.Warnf("failed to unmarshal SubscritpionEvent:%s", err.Error())
			return
		}
		h.sender.SendConfirmation(ctx, mappers.SubscribeEventToDTO(e))

	case ConfirmedRoute:
		e := &events.ConfirmedEvent{}
		if err := proto.Unmarshal(body, e); err != nil {
			h.logger.Warnf("failed to unmarshal ConfirmedEvent:%s", err.Error())
			return
		}
		h.sender.SendConfirmed(ctx, mappers.ConfirmEventToDTO(e))

	case WeatherSuccessRoute:
		e := &events.WeatherSuccessEvent{}
		if err := proto.Unmarshal(body, e); err != nil {
			h.logger.Warnf("failed to unmarshal WeatherSuccessEvent:%s", err.Error())
			return
		}
		h.sender.SendWeather(ctx, mappers.SuccessWeatherToDTO(e))

	case WeatherErrorRoute:
		e := &events.WeatherErrorEvent{}
		if err := proto.Unmarshal(body, e); err != nil {
			h.logger.Warnf("failed to unmarshal WeatherErrorEvent:%s", err.Error())
			return
		}
		h.sender.SendError(ctx, mappers.ErrorWeatherToDTO(e))

	default:
		h.logger.Warnf("unknown event: %s", routingKey)

	}

}
