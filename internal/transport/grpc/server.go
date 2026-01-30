package grpc

import (
	authuc "github.com/SRINIVAS-B-SINDAGI/employee-api/internal/usecase/auth"
	employeeuc "github.com/SRINIVAS-B-SINDAGI/employee-api/internal/usecase/employee"
	authv1 "github.com/SRINIVAS-B-SINDAGI/employee-api/proto/auth/v1"
	employeev1 "github.com/SRINIVAS-B-SINDAGI/employee-api/proto/employee/v1"
	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ServerConfig holds the configuration for the gRPC server.
type ServerConfig struct {
	AuthService     authuc.Service
	EmployeeService employeeuc.Service
	Logger          log.Logger
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
	authv1.RegisterAuthServiceServer(server, NewAuthServer(cfg.AuthService))
	employeev1.RegisterEmployeeServiceServer(server, NewEmployeeServer(cfg.EmployeeService))

	// Enable reflection for grpcurl and other tools
	reflection.Register(server)

	return server
}
