package grpc

import (
	"fmt"
	"time"
	"weather-forecast/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectWithRetry(address string, cfg Config, log logger.Logger) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error

	for i := 0; i < cfg.Retries; i++ {
		log.Infof("Attempting to connect to gRPC service %s (attempt %d/%d)...", address, i+1, cfg.Retries)

		conn, err = grpc.NewClient(address,
			grpc.WithUnaryInterceptor(CorrelationIDClientInterceptor()),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err == nil {
			log.Infof("Successfully connected to gRPC service: %s", address)
			return conn, nil
		}

		log.Warnf("Failed to connect to gRPC service %s: %v", address, err)

		if i < cfg.Retries-1 {
			log.Infof("Retrying in %d seconds...", cfg.RetryDelay)
			time.Sleep(time.Duration(cfg.RetryDelay) * time.Second)
		}
	}

	return nil, fmt.Errorf("failed to connect to gRPC service %s after %d attempts: %w", address, cfg.Retries, err)
}
