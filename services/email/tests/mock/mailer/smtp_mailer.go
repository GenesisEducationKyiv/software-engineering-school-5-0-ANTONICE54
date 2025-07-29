package mailer

import (
	"context"
	"sync"
)

type (
	SentEmail struct {
		Subject string
		Body    string
		SentTo  string
	}
	MockSMTPMailer struct {
		sentEmails []SentEmail
		mu         *sync.RWMutex
	}
)

func NewMockSMTPMailer() *MockSMTPMailer {
	return &MockSMTPMailer{
		sentEmails: make([]SentEmail, 0),
		mu:         &sync.RWMutex{},
	}
}

func (m *MockSMTPMailer) Send(ctx context.Context, subject string, body, email string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sentEmails = append(m.sentEmails, SentEmail{
		Subject: subject,
		Body:    body,
		SentTo:  email,
	})

	return nil
}

func (m *MockSMTPMailer) GetSentEmails() []SentEmail {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]SentEmail, len(m.sentEmails))
	copy(result, m.sentEmails)
	return result
}
