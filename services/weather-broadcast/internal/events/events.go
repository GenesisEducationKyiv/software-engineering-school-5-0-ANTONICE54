package events

import (
	"errors"
	"weather-broadcast-service/internal/dto"
	protoevents "weather-forecast/pkg/proto/events"

	"google.golang.org/protobuf/proto"
)

const (
	weatherSuccessRoute = "emails.weather.success"
	weatherErrorRoute   = "emails.weather.error"
)

type (
	EventType string

	Event struct {
		Type EventType
		Body []byte
	}
)

var (
	weatherSuccessEvent EventType = "Weather success"
	weatherErrorEvent   EventType = "Weather error"
)

func NewWeatherSuccess(info *dto.WeatherMailSuccessInfo) (*Event, error) {
	e := &protoevents.WeatherSuccessEvent{
		Email: info.Email,
		City:  info.City,
		Weather: &protoevents.Weather{
			Temperature: info.Weather.Temperature,
			Humidity:    int32(info.Weather.Humidity),
			Description: info.Weather.Description,
		},
	}
	body, err := proto.Marshal(e)

	if err != nil {
		return nil, err
	}

	return &Event{
		Type: weatherSuccessEvent,
		Body: body,
	}, nil
}

func NewWeatherError(info *dto.WeatherMailErrorInfo) (*Event, error) {
	e := &protoevents.WeatherErrorEvent{
		Email: info.Email,
		City:  info.City,
	}
	body, err := proto.Marshal(e)

	if err != nil {
		return nil, err
	}

	return &Event{
		Type: weatherErrorEvent,
		Body: body,
	}, nil
}

func (e *Event) RoutingKey() (string, error) {

	switch e.Type {
	case weatherSuccessEvent:
		return weatherSuccessRoute, nil
	case weatherErrorEvent:
		return weatherErrorRoute, nil
	default:
		return "", errors.New("unknown route key")
	}

}
