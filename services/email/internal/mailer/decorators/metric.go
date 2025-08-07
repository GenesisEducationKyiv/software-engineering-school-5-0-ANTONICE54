package decorators

import (
	"context"
	"email-service/internal/metrics"
	"email-service/internal/services"
	"weather-forecast/pkg/logger"
)

type (
	MetricMailer struct {
		mailer services.Mailer
		metric metrics.MailRecorder
		logger logger.Logger
	}
)

func NewMetricMailer(mailer services.Mailer, metric metrics.MailRecorder, logger logger.Logger) *MetricMailer {
	return &MetricMailer{
		mailer: mailer,
		metric: metric,
		logger: logger,
	}
}

func (m *MetricMailer) Send(ctx context.Context, subject string, body, email string) error {
	log := m.logger.WithContext(ctx)

	err := m.mailer.Send(ctx, subject, body, email)
	if err != nil {
		log.Debugf("Incrementing emails_failed_total metric")
		m.metric.RecordEmailFail(subject)
		return err
	}

	log.Debugf("Incrementing emails_sent_total metric")
	m.metric.RecordEmailSuccess(subject)

	return nil

}
