package grpc

import (
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/infrastructure/auth"
	authuc "github.com/SRINIVAS-B-SINDAGI/employee-api/internal/usecase/auth"
	employeeuc "github.com/SRINIVAS-B-SINDAGI/employee-api/internal/usecase/employee"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/usecase/salary"
	authv1 "github.com/SRINIVAS-B-SINDAGI/employee-api/proto/auth/v1"
	employeev1 "github.com/SRINIVAS-B-SINDAGI/employee-api/proto/employee/v1"
	salaryv1 "github.com/SRINIVAS-B-SINDAGI/employee-api/proto/salary/v1"
	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ServerConfig holds the configuration for the gRPC server.
type ServerConfig struct {
	AuthService     authuc.Service
	EmployeeService employeeuc.Service
	SalaryService   salary.Service
	Logger          log.Logger
	JWTManager      *auth.JWTManager
}

// NewServer creates a new gRPC server with all services registered.
func NewServer(cfg ServerConfig) *grpc.Server {
	// Create interceptor chain
	interceptors := grpc.ChainUnaryInterceptor(
		RecoveryInterceptor(cfg.Logger),
		LoggingInterceptor(cfg.Logger),
		AuthInterceptor(cfg.JWTManager),
	)

	// Create gRPC server with interceptors
	server := grpc.NewServer(interceptors)
	authv1.RegisterAuthServiceServer(server, NewAuthServer(cfg.AuthService))
	employeev1.RegisterEmployeeServiceServer(server, NewEmployeeServer(cfg.EmployeeService))
	salaryv1.RegisterSalaryServiceServer(server, NewSalaryServer(cfg.SalaryService))
	// Enable reflection for grpcurl and other tools
	reflection.Register(server)

	return server
}
