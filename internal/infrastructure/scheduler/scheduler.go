package scheduler

import (
	"time"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/logger"

	"github.com/robfig/cron/v3"
)

const DAILY = "0 12 * * * "
const HOURLY = "9 * * * * "

type (
	NotificationServiceI interface {
		SendWeather(frequency models.Frequency)
	}

	Scheduler struct {
		cron                cron.Cron
		notificationService NotificationServiceI
		logger              logger.Logger
	}
)

func New(notificationService NotificationServiceI, location *time.Location, logger logger.Logger) *Scheduler {
	return &Scheduler{
		cron:                *cron.New(cron.WithLocation(location)),
		notificationService: notificationService,
		logger:              logger,
	}

}

func (s *Scheduler) Init() {
	_, err := s.cron.AddFunc(DAILY, func() { s.notificationService.SendWeather(models.Daily) })
	if err != nil {
		s.logger.Fatalf("Failed to setup daily sender: %s", err.Error())
		return
	}
	_, err = s.cron.AddFunc(HOURLY, func() { s.notificationService.SendWeather(models.Hourly) })
	if err != nil {
		s.logger.Fatalf("Failed to setup hourly sender: %s", err.Error())
		return
	}
}

func (s *Scheduler) Run() {
	s.cron.Start()
}
