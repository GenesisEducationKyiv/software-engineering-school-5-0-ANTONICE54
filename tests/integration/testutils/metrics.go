package testutils

import "sync"

type InMemoryMetrics struct {
	CacheHit   int
	CacheMiss  int
	CacheError int
	mu         *sync.Mutex
}

func NewInMemoryMetrics() *InMemoryMetrics {
	return &InMemoryMetrics{
		CacheHit:   0,
		CacheMiss:  0,
		CacheError: 0,
		mu:         &sync.Mutex{},
	}
}

func (m *InMemoryMetrics) RecordCacheHit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CacheHit++

}

func (m *InMemoryMetrics) RecordCacheMiss() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CacheMiss++
}

func (m *InMemoryMetrics) RecordCacheError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CacheError++
}

func (m *InMemoryMetrics) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CacheHit = 0
	m.CacheMiss = 0
	m.CacheError = 0
}

func (m *InMemoryMetrics) Stats() (int, int, int) {
	return m.CacheHit, m.CacheMiss, m.CacheError
}
