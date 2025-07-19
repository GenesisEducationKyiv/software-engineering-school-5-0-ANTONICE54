package publisher

import (
	"context"
	"encoding/json"
	infraerror "subscription-service/internal/infrastructure/errors"

	"sync"
	"weather-forecast/pkg/events"
)

type PublishedEvent struct {
	EventType events.EventType
	RawData   []byte
}

type MockEventPublisher struct {
	publishedEvents []PublishedEvent
	mu              sync.RWMutex
	shouldFail      bool
}

func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{
		publishedEvents: make([]PublishedEvent, 0),
	}
}

func (m *MockEventPublisher) Publish(ctx context.Context, event events.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	body, err := json.Marshal(event)
	if err != nil {
		return infraerror.InternalError
	}

	m.publishedEvents = append(m.publishedEvents, PublishedEvent{
		EventType: event.EventType(),
		RawData:   body,
	})

	return nil
}

func (m *MockEventPublisher) GetPublishedEvents() []PublishedEvent {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]PublishedEvent, len(m.publishedEvents))
	copy(result, m.publishedEvents)
	return result
}
