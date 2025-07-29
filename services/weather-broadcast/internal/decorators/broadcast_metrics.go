package decorators

import (
	"context"
	"time"
	"weather-broadcast-service/internal/metrics"
	"weather-broadcast-service/internal/models"
	"weather-broadcast-service/internal/scheduler"
	"weather-forecast/pkg/logger"
)

type BroadcastMetricsDecorator struct {
	service scheduler.WeatherBroadcastService
	metrics metrics.BroadcastRecorder
	logger  logger.Logger
}

func NewBroadcastMetricsDecorator(service scheduler.WeatherBroadcastService, metrics metrics.BroadcastRecorder, logger logger.Logger) *BroadcastMetricsDecorator {
	return &BroadcastMetricsDecorator{
		service: service,
		metrics: metrics,
		logger:  logger,
	}
}

func (d *BroadcastMetricsDecorator) Broadcast(ctx context.Context, frequency models.Frequency) {
	log := d.logger.WithContext(ctx)

	start := time.Now()

	log.Infof("Starting weather broadcast for frequency: %s", frequency)

	d.service.Broadcast(ctx, frequency)

	duration := time.Since(start)
	d.metrics.RecordBroadcastDuration(string(frequency), duration)

	log.Infof("Weather broadcast completed for %s subscription in %v", frequency, duration)
}
