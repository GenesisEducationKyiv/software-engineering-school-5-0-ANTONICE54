package server

import (
	"net"
	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/proto/weather"

	"google.golang.org/grpc"
)

type (
	Server struct {
		grpcServer *grpc.Server
		logger     logger.Logger
	}
)

func New(weatherHandler weather.WeatherServiceServer, logger logger.Logger) *Server {
	grpcServer := grpc.NewServer()

	weather.RegisterWeatherServiceServer(grpcServer, weatherHandler)

	return &Server{
		grpcServer: grpcServer,
		logger:     logger,
	}
}

func (s *Server) Start(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s.logger.Infof("gRPC server starting at port %s:", port)
	return s.grpcServer.Serve(lis)
}
