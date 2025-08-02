package logger

import "sync"

type (
	Sampler interface {
		ShouldLog() bool
	}
	RateSampler struct {
		rate    int64
		counter int64
		mu      sync.Mutex
	}
	NoSampler struct {
	}
)

func NewRateSampler(rate int) *RateSampler {
	return &RateSampler{rate: int64(rate)}
}

func (s *RateSampler) ShouldLog() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	return s.counter%s.rate == 0
}

func (s *NoSampler) ShouldLog() bool {
	return true
}
