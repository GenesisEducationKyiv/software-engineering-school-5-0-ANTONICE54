package server

import (
	"weather-forecast/gateway/internal/server/middleware"
	"weather-forecast/pkg/logger"

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

	Server struct {
		router               *gin.Engine
		weatherHandler       WeatherHandler
		subscrtiptionHandler SubscriptionHandler
		logger               logger.Logger
	}
)

func New(weatherHandeler WeatherHandler, subscrtiptionHandler SubscriptionHandler, logger logger.Logger) *Server {

	s := &Server{
		router:               gin.Default(),
		weatherHandler:       weatherHandeler,
		subscrtiptionHandler: subscrtiptionHandler,
		logger:               logger,
	}
	s.setUpRoutes()
	s.setUpMiddleware()
	return s
}

func (s *Server) setUpMiddleware() {
	s.router.Use(middleware.ProcessIDMiddleware())
}

func (s *Server) setUpRoutes() {
	s.router.GET("/", func(ctx *gin.Context) {
		ctx.File("./static/subscription.html")
	})
	s.router.GET("/weather", s.weatherHandler.Get)
	s.router.POST("/subscribe", s.subscrtiptionHandler.Subscribe)
	s.router.GET("/confirm/:token", s.subscrtiptionHandler.Confirm)
	s.router.GET("/unsubscribe/:token", s.subscrtiptionHandler.Unsubscribe)

}

func (s *Server) Run(port string) {
	err := s.router.Run("0.0.0.0:" + port)
	if err != nil {
		s.logger.Fatalf("Failed to start server: %s", err.Error())
	}
}
