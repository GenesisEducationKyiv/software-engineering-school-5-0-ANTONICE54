package scheduler

import (
	"context"
	"time"
	"weather-broadcast-service/internal/models"
	"weather-forecast/pkg/logger"

	"github.com/robfig/cron/v3"
)

const DAILY = "0 12 * * * " //every day at 12 am
const HOURLY = "3 * * * * " //every hour at 0 minute

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

	s.logger.Infof("Setting up scheduler with daily and hourly broadcasts")

	_, err := s.cron.AddFunc(DAILY, func() {
		s.logger.Infof("Daily broadcast triggered")

		s.broadcastService.Broadcast(s.ctx, models.Daily)
	})
	if err != nil {
		s.logger.Fatalf("Failed to setup daily sender: %s", err.Error())
		return
	}
	_, err = s.cron.AddFunc(HOURLY, func() {
		s.logger.Infof("Hourly broadcast triggered")
		s.broadcastService.Broadcast(s.ctx, models.Hourly)
	})
	if err != nil {
		s.logger.Fatalf("Failed to setup hourly sender: %s", err.Error())
		return
	}

	s.logger.Infof("Scheduler setup completed successfully")

}

func (s *Scheduler) Run() {
	s.logger.Infof("Starting scheduler")

	s.cron.Start()

	s.logger.Infof("Scheduler started successfully")

}
