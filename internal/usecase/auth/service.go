package auth

import (
	"context"

	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/entity"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/repository"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/pkg/errors"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, email, password string) (*entity.User, error)
}

type service struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) Service {
	return &service{
		userRepo: userRepo,
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
