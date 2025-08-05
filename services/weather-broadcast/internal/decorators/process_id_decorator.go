package decorators

import (
	"context"
	"weather-broadcast-service/internal/models"
	"weather-broadcast-service/internal/scheduler"
	"weather-forecast/pkg/ctxutil"
	"weather-forecast/pkg/logger"

	"github.com/google/uuid"
)

type CorrelationIDDecorator struct {
	service scheduler.WeatherBroadcastService
	logger  logger.Logger
}

func NewCorrelationIDDecorator(service scheduler.WeatherBroadcastService, logger logger.Logger) *CorrelationIDDecorator {
	return &CorrelationIDDecorator{
		service: service,
		logger:  logger,
	}
}

func (d *CorrelationIDDecorator) Broadcast(ctx context.Context, frequency models.Frequency) {
	correlationID := uuid.New().String()

	//nolint:staticcheck
	ctx = context.WithValue(ctx, ctxutil.CorrelationIDKey.String(), correlationID)
	d.service.Broadcast(ctx, frequency)

}
