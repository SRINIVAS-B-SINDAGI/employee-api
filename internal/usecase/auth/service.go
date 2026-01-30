package auth

import (
	"context"

	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/entity"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/repository"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/infrastructure/auth"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/pkg/errors"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, email, password string) (*entity.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type service struct {
	userRepo   repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewService(userRepo repository.UserRepository, jwtManager *auth.JWTManager) Service {
	return &service{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *service) Register(ctx context.Context, email, password string) (*entity.User, error) {
	if err := validator.ValidateEmail(email); err != nil {
		return nil, err
	}
	if err := validator.ValidatePassword(password); err != nil {
		return nil, err
	}

	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.NewConflictError("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	user := entity.NewUser(email, string(hashedPassword))
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	if err := validator.ValidateEmail(email); err != nil {
		return "", err
	}
	if password == "" {
		return "", errors.NewValidationError("password is required")
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return "", errors.NewUnauthorizedError("invalid credentials")
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.NewUnauthorizedError("invalid credentials")
	}

	token, err := s.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", errors.NewInternalError(err)
	}

	return token, nil
}
