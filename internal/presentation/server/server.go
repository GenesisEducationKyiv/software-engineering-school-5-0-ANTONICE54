package server

import (
	"weather-forecast/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
)

type (
	WeatherHandler interface {
		Get(ctx *gin.Context)
	}
	SubscriptionHandler interface {
		Subscribe(ctx *gin.Context)
		Confirm(ctx *gin.Context)
		Unsubscribe(ctx *gin.Context)
	}

	Scheduler interface {
		SetUp()
		Run()
	}

	Server struct {
		router              *gin.Engine
		subscriptionHandler SubscriptionHandler
		weatherHandler      WeatherHandler
		scheduler           Scheduler
		logger              logger.Logger
	}
)

func New(subscriptionHandler SubscriptionHandler, weatherHandeler WeatherHandler, scheduler Scheduler, logger logger.Logger) *Server {

	s := &Server{
		router:              gin.Default(),
		subscriptionHandler: subscriptionHandler,
		weatherHandler:      weatherHandeler,
		scheduler:           scheduler,
		logger:              logger,
	}
	s.scheduler.SetUp()
	s.setUpRoutes()
	return s
}

func (s *Server) setUpRoutes() {
	s.router.GET("/", func(ctx *gin.Context) {
		ctx.File("./subscription.html")
	})
	s.router.GET("/weather", s.weatherHandler.Get)
	s.router.POST("/subscribe", s.subscriptionHandler.Subscribe)
	s.router.GET("/confirm/:token", s.subscriptionHandler.Confirm)
	s.router.GET("/unsubscribe/:token", s.subscriptionHandler.Unsubscribe)

}

func (s *Server) Run(port string) {
	s.scheduler.Run()
	err := s.router.Run("0.0.0.0:" + port)
	if err != nil {
		s.logger.Fatalf("Failed to start server: %s", err.Error())
	}
}
