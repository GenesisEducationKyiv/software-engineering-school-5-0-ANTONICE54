package decorators

import (
	"context"
	"weather-broadcast-service/internal/models"
	"weather-broadcast-service/internal/scheduler"
	"weather-forecast/pkg/ctxutil"
	"weather-forecast/pkg/logger"

	"github.com/google/uuid"
)

const processIDKey = "process_id"

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

	ctx = context.WithValue(ctx, ctxutil.ProcessIDKey, processID)
	d.service.Broadcast(ctx, frequency)

}
