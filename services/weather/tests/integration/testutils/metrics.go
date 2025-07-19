package testutils

import "sync"

type InMemoryMetrics struct {
	cacheHit   int
	cacheMiss  int
	cacheError int
	mu         *sync.Mutex
}

func NewInMemoryMetrics() *InMemoryMetrics {
	return &InMemoryMetrics{
		cacheHit:   0,
		cacheMiss:  0,
		cacheError: 0,
		mu:         &sync.Mutex{},
	}
}

func (m *InMemoryMetrics) RecordCacheHit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheHit++

}

func (m *InMemoryMetrics) RecordCacheMiss() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheMiss++
}

func (m *InMemoryMetrics) RecordCacheError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheError++
}

func (m *InMemoryMetrics) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheHit = 0
	m.cacheMiss = 0
	m.cacheError = 0
}

func (m *InMemoryMetrics) Stats() (int, int, int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.cacheHit, m.cacheMiss, m.cacheError
}
