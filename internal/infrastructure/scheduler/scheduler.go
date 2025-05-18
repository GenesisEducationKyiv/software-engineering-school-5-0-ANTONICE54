package scheduler

import (
	"context"
	"time"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/logger"

	"github.com/robfig/cron/v3"
)

const DAILY = "0 12 * * * "
const HOURLY = "16 * * * * "

type (
	NotificationService interface {
		SendWeather(ctx context.Context, frequency models.Frequency)
	}

	Scheduler struct {
		cron                cron.Cron
		notificationService NotificationService
		logger              logger.Logger
	}
)

func New(notificationService NotificationService, location *time.Location, logger logger.Logger) *Scheduler {
	return &Scheduler{
		cron:                *cron.New(cron.WithLocation(location)),
		notificationService: notificationService,
		logger:              logger,
	}

}

func (s *Scheduler) SetUp() {

	_, err := s.cron.AddFunc(DAILY, func() { s.notificationService.SendWeather(context.Background(), models.Daily) })
	if err != nil {
		s.logger.Fatalf("Failed to setup daily sender: %s", err.Error())
		return
	}
	_, err = s.cron.AddFunc(HOURLY, func() { s.notificationService.SendWeather(context.Background(), models.Hourly) })
	if err != nil {
		s.logger.Fatalf("Failed to setup hourly sender: %s", err.Error())
		return
	}
}

func (s *Scheduler) Run() {
	s.cron.Start()
}
