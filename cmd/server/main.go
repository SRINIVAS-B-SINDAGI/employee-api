package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/infrastructure/auth"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/infrastructure/config"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/infrastructure/persistence/postgres"
	transportgrpc "github.com/SRINIVAS-B-SINDAGI/employee-api/internal/transport/grpc"
	authuc "github.com/SRINIVAS-B-SINDAGI/employee-api/internal/usecase/auth"
	employeeuc "github.com/SRINIVAS-B-SINDAGI/employee-api/internal/usecase/employee"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func main() {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	cfg, err := config.Load()
	if err != nil {
		_ = level.Error(logger).Log("msg", "failed to load config", "err", err)
		os.Exit(1)
	}

	_ = level.Info(logger).Log("msg", "connecting to database")
	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		_ = level.Error(logger).Log("msg", "failed to connect to database", "err", err)
		os.Exit(1)
	}

	defer func() {
		if err := postgres.Close(db); err != nil {
			_ = level.Error(logger).Log("msg", "failed to close database", "err", err)
		}
	}()

	_ = level.Info(logger).Log("msg", "running database migrations")
	if err := postgres.AutoMigrate(db); err != nil {
		_ = level.Error(logger).Log("msg", "failed to run migrations", "err", err)
		os.Exit(1)
	}

	userRepo := postgres.NewUserRepository(db)
	employeeRepo := postgres.NewEmployeeRepository(db)
	jwtManager := auth.NewJWTManager(cfg.JWT)
	authService := authuc.NewService(userRepo, jwtManager)

	employeeService := employeeuc.NewService(employeeRepo)
	grpcServer := transportgrpc.NewServer(transportgrpc.ServerConfig{
		AuthService:     authService,
		EmployeeService: employeeService,
		Logger:          log.With(logger, "transport", "grpc"),
	})

	errChan := make(chan error, 1)
	go func() {
		grpcListener, err := net.Listen("tcp", ":"+cfg.Server.GRPCPort)
		if err != nil {
			errChan <- err
			return
		}
		_ = level.Info(logger).Log("msg", "starting gRPC server", "port", cfg.Server.GRPCPort)
		errChan <- grpcServer.Serve(grpcListener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		_ = level.Error(logger).Log("msg", "server error", "err", err)
	case sig := <-quit:
		_ = level.Info(logger).Log("msg", "shutting down servers", "signal", sig)
	}

	grpcServer.GracefulStop()

	_ = level.Info(logger).Log("msg", "servers stopped")
}
