package publisher

import (
	"context"

	"sync"
)

type PublishedEvent struct {
	EventType string
	RawData   []byte
}

type MockEventPublisher struct {
	publishedEvents []PublishedEvent
	mu              sync.RWMutex
}

func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{
		publishedEvents: make([]PublishedEvent, 0),
	}
}

func (m *MockEventPublisher) Publish(ctx context.Context, routingKey string, body []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.publishedEvents = append(m.publishedEvents, PublishedEvent{
		EventType: routingKey,
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
