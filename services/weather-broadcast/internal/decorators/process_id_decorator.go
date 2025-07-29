package decorators

import (
	"context"
	"weather-broadcast-service/internal/models"
	"weather-broadcast-service/internal/scheduler"
	"weather-forecast/pkg/logger"

	"github.com/google/uuid"
)

type ProcessIDDecorator struct {
	service scheduler.WeatherBroadcastService
	logger  logger.Logger
}

func NewProcessIDDecorator(service scheduler.WeatherBroadcastService, logger logger.Logger) *ProcessIDDecorator {
	return &ProcessIDDecorator{
		service: service,
		logger:  logger,
	}
}

func (d *ProcessIDDecorator) Broadcast(ctx context.Context, frequency models.Frequency) {
	processID := uuid.New().String()
	ctx = context.WithValue(ctx, "process_id", processID)
	d.service.Broadcast(ctx, frequency)

}
