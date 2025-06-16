package stub_services

import (
	"context"
	"weather-forecast/internal/domain/models"
)

type StubMailer struct{}

func NewStubMailer() *StubMailer {
	return &StubMailer{}
}

func (m *StubMailer) SendConfirmation(ctx context.Context, email, token string, frequency models.Frequency) {
}
func (m *StubMailer) SendConfirmed(ctx context.Context, email, token string, frequency models.Frequency) {
}
