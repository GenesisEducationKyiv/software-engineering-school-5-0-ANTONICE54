package scheduler

import (
	"context"
	"time"
	"weather-broadcast-service/internal/models"
	"weather-forecast/pkg/logger"

	"github.com/robfig/cron/v3"
)

const DAILY = "00 12 * * * " //every day at 12 am
const HOURLY = "0 * * * * "  //every hour at 0 minute

type (
	WeatherBroadcastService interface {
		Broadcast(ctx context.Context, frequency models.Frequency)
	}

	Scheduler struct {
		cron             cron.Cron
		broadcastService WeatherBroadcastService
		logger           logger.Logger
		ctx              context.Context
	}
)

func New(ctx context.Context, notificationService WeatherBroadcastService, location *time.Location, logger logger.Logger) *Scheduler {
	return &Scheduler{
		cron:             *cron.New(cron.WithLocation(location)),
		broadcastService: notificationService,
		logger:           logger,
		ctx:              ctx,
	}

}

func (s *Scheduler) SetUp() {

	_, err := s.cron.AddFunc(DAILY, func() { s.broadcastService.Broadcast(s.ctx, models.Daily) })
	if err != nil {
		s.logger.Fatalf("Failed to setup daily sender: %s", err.Error())
		return
	}
	_, err = s.cron.AddFunc(HOURLY, func() {

		s.broadcastService.Broadcast(s.ctx, models.Hourly)
	})
	if err != nil {
		s.logger.Fatalf("Failed to setup hourly sender: %s", err.Error())
		return
	}
}

func (s *Scheduler) Run() {
	s.cron.Start()
}
