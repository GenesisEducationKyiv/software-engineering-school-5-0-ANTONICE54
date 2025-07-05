package stub_services

import (
	"context"
	"weather-forecast/internal/domain/models"
)

type (
	SentConfirmation struct {
		Email     string
		Frequency models.Frequency
	}

	SentConfirmed struct {
		Email     string
		Frequency models.Frequency
	}

	StubMailer struct {
		SentConfirmations []SentConfirmation
		SentConfirmeds    []SentConfirmed
	}
)

func NewStubMailer() *StubMailer {
	return &StubMailer{}
}

func (m *StubMailer) SendConfirmation(ctx context.Context, email, token string, frequency models.Frequency) {
	content := SentConfirmation{
		Email:     email,
		Frequency: frequency,
	}
	m.SentConfirmations = append(m.SentConfirmations, content)
}

func (m *StubMailer) SendConfirmed(ctx context.Context, email, token string, frequency models.Frequency) {
	content := SentConfirmed{
		Email:     email,
		Frequency: frequency,
	}
	m.SentConfirmeds = append(m.SentConfirmeds, content)
}
