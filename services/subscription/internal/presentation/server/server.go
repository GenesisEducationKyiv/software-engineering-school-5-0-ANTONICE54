package server

import (
	"net"
	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/proto/subscription"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type (
	Server struct {
		grpcServer *grpc.Server
		logger     logger.Logger
	}
)

func New(subscriptionHandler subscription.SubscriptionServiceServer, logger logger.Logger) *Server {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	subscription.RegisterSubscriptionServiceServer(grpcServer, subscriptionHandler)

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
