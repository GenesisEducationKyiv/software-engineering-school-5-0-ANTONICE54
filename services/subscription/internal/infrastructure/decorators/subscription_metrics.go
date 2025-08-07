package decorators

import (
	"context"
	"subscription-service/internal/domain/models"
	"subscription-service/internal/domain/usecases"
	"subscription-service/internal/infrastructure/metrics"
	"weather-forecast/pkg/logger"
)

type SubscriptionServiceMetricsDecorator struct {
	service usecases.SubscriptionService
	metrics metrics.SubscriptionRecorder
	logger  logger.Logger
}

func NewSubscriptionServiceMetricsDecorator(service usecases.SubscriptionService, metrics metrics.SubscriptionRecorder, logger logger.Logger) *SubscriptionServiceMetricsDecorator {
	return &SubscriptionServiceMetricsDecorator{
		service: service,
		metrics: metrics,
		logger:  logger,
	}
}

func (d *SubscriptionServiceMetricsDecorator) Subscribe(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
	log := d.logger.WithContext(ctx)
	result, err := d.service.Subscribe(ctx, subscription)

	if err == nil && result != nil {
		log.Debugf("Incrementing subscriptions_created_total metric")
		d.metrics.RecordSubscriptionCreated()
	}

	return result, err
}

func (d *SubscriptionServiceMetricsDecorator) Confirm(ctx context.Context, token string) error {
	log := d.logger.WithContext(ctx)

	err := d.service.Confirm(ctx, token)

	if err == nil {
		log.Debugf("Incrementing subscriptions_confirmed_total metric")
		d.metrics.RecordSubscriptionConfirmed()
	}

	return err
}

func (d *SubscriptionServiceMetricsDecorator) Unsubscribe(ctx context.Context, token string) error {
	log := d.logger.WithContext(ctx)

	err := d.service.Unsubscribe(ctx, token)

	if err == nil {
		log.Debugf("Incrementing subscriptions_deleted_total metric")
		d.metrics.RecordSubscriptionDeleted()
	}

	return err
}

func (d *SubscriptionServiceMetricsDecorator) ListByFrequency(ctx context.Context, query *models.ListSubscriptionsQuery) ([]models.Subscription, error) {
	return d.service.ListByFrequency(ctx, query)
}
