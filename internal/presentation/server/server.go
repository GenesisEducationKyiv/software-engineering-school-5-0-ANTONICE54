package server

import (
	"weather-forecast/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
)

type (
	WeatherHandlerI interface {
		Get(ctx *gin.Context)
	}
	SubscriptionHandlerI interface {
		Subscribe(ctx *gin.Context)
		Confirm(ctx *gin.Context)
		Unsubscribe(ctx *gin.Context)
	}

	SchedulerI interface {
		Init()
		Run()
	}

	Server struct {
		router              *gin.Engine
		subscriptionHandler SubscriptionHandlerI
		weatherHandler      WeatherHandlerI
		scheduler           SchedulerI
		logger              logger.Logger
	}
)

func New(subscriptionHandler SubscriptionHandlerI, weatherHandeler WeatherHandlerI, scheduler SchedulerI, logger logger.Logger) *Server {

	s := &Server{
		router:              gin.Default(),
		subscriptionHandler: subscriptionHandler,
		weatherHandler:      weatherHandeler,
		scheduler:           scheduler,
		logger:              logger,
	}
	s.scheduler.Init()
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
