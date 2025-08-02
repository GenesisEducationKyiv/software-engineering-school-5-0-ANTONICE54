package events

import (
	"subscription-service/internal/domain/contracts"
	infraerror "subscription-service/internal/infrastructure/errors"
	protoevents "weather-forecast/pkg/proto/events"

	"google.golang.org/protobuf/proto"
)

const (
	confirmationRoute = "emails.subscription"
	confirmedRoute    = "emails.confirmed"
	unsubscribedRoute = "emails.unsubscribed"

	confirmedEvent    EventType = "CONFIRMED"
	unsubscribedEvent EventType = "UNSUBSCRIBED"
	confirmationEvent EventType = "CONFIRMATION"
)

type (
	EventType string

	Event struct {
		Type EventType
		Body []byte
	}
)

func NewConfirmation(info *contracts.ConfirmationInfo) (*Event, error) {
	e := &protoevents.SubscriptionEvent{
		Email:     info.Email,
		Token:     info.Token,
		Frequency: string(info.Frequency),
	}

	body, err := proto.Marshal(e)
	if err != nil {
		return nil, err
	}

	return &Event{
		Type: confirmationEvent,
		Body: body,
	}, nil
}

func NewConfirmed(info *contracts.ConfirmedInfo) (*Event, error) {
	e := &protoevents.ConfirmedEvent{
		Email:     info.Email,
		Token:     info.Token,
		Frequency: string(info.Frequency),
	}

	body, err := proto.Marshal(e)
	if err != nil {
		return nil, err
	}

	return &Event{
		Type: confirmedEvent,
		Body: body,
	}, nil
}

func NewUnsubscribed(info *contracts.UnsubscribeInfo) (*Event, error) {
	e := &protoevents.UnsubscribedEvent{
		Email:     info.Email,
		City:      info.City,
		Frequency: string(info.Frequency),
	}

	body, err := proto.Marshal(e)
	if err != nil {
		return nil, err
	}

	return &Event{
		Type: unsubscribedEvent,
		Body: body,
	}, nil
}

func (e *Event) RoutingKey() (string, error) {

	switch e.Type {
	case confirmationEvent:
		return confirmationRoute, nil
	case confirmedEvent:
		return confirmedRoute, nil
	case unsubscribedEvent:
		return unsubscribedRoute, nil
	default:
		return "", infraerror.ErrUnknownEventRoute
	}

}
