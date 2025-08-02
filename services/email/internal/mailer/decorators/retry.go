package decorators

import (
	"context"
	"email-service/internal/config"
	"email-service/internal/services"

	"time"
	"weather-forecast/pkg/logger"
)

type (
	RetryMailer struct {
		mailer     services.Mailer
		maxRetries int
		delay      time.Duration
		logger     logger.Logger
	}
)

func NewRetryMailer(mailer services.Mailer, cfg config.Retry, logger logger.Logger) *RetryMailer {

	duration := time.Duration(cfg.Delay) * time.Second

	return &RetryMailer{
		mailer:     mailer,
		maxRetries: cfg.MaxRetries,
		delay:      duration,
		logger:     logger,
	}

}

func (m *RetryMailer) Send(ctx context.Context, subject string, body, email string) error {
	log := m.logger.WithContext(ctx)

	var err error

	for attempt := 0; attempt < m.maxRetries; attempt++ {

		err = m.mailer.Send(ctx, subject, body, email)

		if err != nil {
			if attempt < m.maxRetries-1 {
				log.Warnf("Attempt %d failed for email to %s, retrying in %v. Error: %s", attempt+1, email, m.delay, err.Error())

				time.Sleep(m.delay)
			} else {
				log.Errorf("Final attempt %d failed for email to %s. Error: %s", attempt+1, email, err.Error())
			}

			continue
		}

		return nil
	}
	return err
}
