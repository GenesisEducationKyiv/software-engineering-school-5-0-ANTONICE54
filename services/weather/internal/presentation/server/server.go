package server

import (
	"net"
	grpcpkg "weather-forecast/pkg/grpc"
	"weather-forecast/pkg/logger"

	"weather-forecast/pkg/proto/weather"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type (
	Server struct {
		grpcServer *grpc.Server
		logger     logger.Logger
	}
)

func New(weatherHandler weather.WeatherServiceServer, logger logger.Logger) *Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcpkg.CorrelationIDServerInterceptor(logger)),
	)
	reflection.Register(grpcServer)

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

func (s *Server) Shutdown() {
	s.grpcServer.GracefulStop()
}
