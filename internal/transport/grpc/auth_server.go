package grpc

import (
	"context"

	authuc "github.com/SRINIVAS-B-SINDAGI/employee-api/internal/usecase/auth"
	authv1 "github.com/SRINIVAS-B-SINDAGI/employee-api/proto/auth/v1"
)

// authServer implements the AuthServiceServer interface.
type authServer struct {
	authv1.UnimplementedAuthServiceServer
	service authuc.Service
}

// NewAuthServer creates a new gRPC auth server.
func NewAuthServer(service authuc.Service) authv1.AuthServiceServer {
	return &authServer{
		service: service,
	}
}

// Register creates a new user account.
func (s *authServer) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	user, err := s.service.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &authv1.RegisterResponse{
		Id:    user.ID.String(),
		Email: user.Email,
	}, nil
}
