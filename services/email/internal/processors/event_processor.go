package processors

import (
	"context"
	"email-service/internal/dto"
	"email-service/internal/mappers"
	"encoding/json"
	"weather-forecast/pkg/events"
	"weather-forecast/pkg/logger"
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

	switch events.EventType(routingKey) {

	case events.SubsctiptionEmail:
		var e events.SubscriptionEvent
		if err := json.Unmarshal(body, &e); err != nil {
			h.logger.Warnf("failed to unmarshal SubscritpionEvent:%s", err.Error())
			return
		}
		h.sender.SendConfirmation(ctx, mappers.SubscribeEventToDTO(e))

	case events.ConfirmedEmail:
		var e events.ConfirmedEvent
		if err := json.Unmarshal(body, &e); err != nil {
			h.logger.Warnf("failed to unmarshal ConfirmedEvent:%s", err.Error())
			return
		}
		h.sender.SendConfirmed(ctx, mappers.ConfirmEventToDTO(e))

	case events.WeatherEmailSuccess:
		var e events.WeatherSuccessEvent
		if err := json.Unmarshal(body, &e); err != nil {
			h.logger.Warnf("failed to unmarshal WeatherSuccessEvent:%s", err.Error())
			return
		}
		h.sender.SendWeather(ctx, mappers.SuccessWeatehrToDTO(e))

	case events.WeatherEmailError:
		var e events.WeatherErrorEvent
		if err := json.Unmarshal(body, &e); err != nil {
			h.logger.Warnf("failed to unmarshal WeatherErrorEvent:%s", err.Error())
			return
		}
		h.sender.SendError(ctx, mappers.ErrorWeatehrToDTO(e))

	default:
		h.logger.Warnf("unknown event: %s", routingKey)

	}

}
