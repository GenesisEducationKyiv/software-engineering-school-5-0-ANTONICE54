package decorators

import (
	"context"
	"time"
	"weather-broadcast-service/internal/metrics"
	"weather-broadcast-service/internal/models"
	"weather-broadcast-service/internal/services"
	"weather-forecast/pkg/logger"
)

type BroadcastMetricsDecorator struct {
	service services.WeatherBroadcastService
	metrics metrics.BroadcastRecorder
	logger  logger.Logger
}

func NewBroadcastMetricsDecorator(service services.WeatherBroadcastService, metrics metrics.BroadcastRecorder, logger logger.Logger) *BroadcastMetricsDecorator {
	return &BroadcastMetricsDecorator{
		service: service,
		metrics: metrics,
		logger:  logger,
	}
}

func (d *BroadcastMetricsDecorator) Broadcast(ctx context.Context, frequency models.Frequency) {
	start := time.Now()

	d.logger.Infof("Starting weather broadcast for frequency: %s", frequency)

	d.service.Broadcast(ctx, frequency)

	duration := time.Since(start)
	d.metrics.RecordBroadcastDuration(string(frequency), duration)

	d.logger.Infof("Weather broadcast completed for %s subscription in %v", frequency, duration)
}
