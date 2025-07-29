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
		SendUnsubscribed(ctx context.Context, info *dto.UnsubscribedEmailInfo)
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
	log := h.logger.WithContext(ctx)

	log.Debugf("Processing event: routing_key=%s, size=%d bytes", routingKey, len(body))

	switch events.EventType(routingKey) {

	case events.SubsctiptionEmail:
		var e events.SubscriptionEvent
		if err := json.Unmarshal(body, &e); err != nil {
			log.Warnf("failed to unmarshal SubscritpionEvent from routing_key = %s:%s", routingKey, err.Error())
			return
		}
		log.Debugf("Successfully parsed SubscriptionEvent for email: %s", e.Email)
		h.sender.SendConfirmation(ctx, mappers.SubscribeEventToDTO(e))

	case events.ConfirmedEmail:
		var e events.ConfirmedEvent
		if err := json.Unmarshal(body, &e); err != nil {
			log.Warnf("failed to unmarshal ConfirmedEvent from routing_key = %s:%s", routingKey, err.Error())
			return
		}
		log.Debugf("Successfully parsed ConfirmedEvent for email: %s", e.Email)

		h.sender.SendConfirmed(ctx, mappers.ConfirmEventToDTO(e))

	case events.UnsubscribedEmail:
		var e events.UnsubscribedEvent
		if err := json.Unmarshal(body, &e); err != nil {
			log.Warnf("failed to unmarshal UnsubscribeEvent from routing_key = %s:%s", routingKey, err.Error())
			return
		}
		log.Debugf("Successfully parsed UnsubscribedEvent for email: %s", e.Email)
		h.sender.SendUnsubscribed(ctx, mappers.UnsubscribeEventToDTO(e))

	case events.WeatherEmailSuccess:
		var e events.WeatherSuccessEvent
		if err := json.Unmarshal(body, &e); err != nil {
			log.Warnf("failed to unmarshal WeatherSuccessEvent from routing_key = %s:%s", routingKey, err.Error())
			return
		}
		log.Debugf("Successfully parsed WeatherSuccessEvent for email: %s, city: %s", e.Email, e.City)
		h.sender.SendWeather(ctx, mappers.SuccessWeatehrToDTO(e))

	case events.WeatherEmailError:
		var e events.WeatherErrorEvent
		if err := json.Unmarshal(body, &e); err != nil {
			log.Warnf("failed to unmarshal WeatherErrorEvent from routing_key = %s:%s", routingKey, err.Error())
			return
		}
		log.Debugf("Successfully parsed WeatherErrorEvent for email: %s, city: %s", e.Email, e.City)
		h.sender.SendError(ctx, mappers.ErrorWeatehrToDTO(e))

	default:
		log.Warnf("unknown event: %s", routingKey)

	}

	log.Debugf("Event processing completed for routing_key: %s", routingKey)

}
