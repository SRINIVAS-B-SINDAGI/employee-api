package grpc

import (
	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ServerConfig holds the configuration for the gRPC server.
type ServerConfig struct {
	Logger log.Logger
}

// NewServer creates a new gRPC server with all services registered.
func NewServer(cfg ServerConfig) *grpc.Server {
	// Create interceptor chain
	interceptors := grpc.ChainUnaryInterceptor(
		RecoveryInterceptor(cfg.Logger),
		LoggingInterceptor(cfg.Logger),
	)

	// Create gRPC server with interceptors
	server := grpc.NewServer(interceptors)

	// Enable reflection for grpcurl and other tools
	reflection.Register(server)

	return server
}
