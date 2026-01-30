package auth

import (
	"context"

	"testing"

	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/entity"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/infrastructure/auth"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/infrastructure/config"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func TestAuthService_Register(t *testing.T) {
	ctx := context.Background()
	jwtManager := auth.NewJWTManager(config.JWTConfig{
		Secret:     "test-secret",
		Expiration: 24 * 60 * 60,
		Issuer:     "test",
	})

	t.Run("successful registration", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := NewService(mockRepo, jwtManager)

		mockRepo.On("ExistsByEmail", ctx, "test@example.com").Return(false, nil)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

		user, err := svc.Register(ctx, "test@example.com", "password123")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test@example.com", user.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("email already exists", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := NewService(mockRepo, jwtManager)

		mockRepo.On("ExistsByEmail", ctx, "existing@example.com").Return(true, nil)

		user, err := svc.Register(ctx, "existing@example.com", "password123")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "already registered")
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid email", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := NewService(mockRepo, jwtManager)

		user, err := svc.Register(ctx, "invalid-email", "password123")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.True(t, errors.IsValidationError(err))
	})

	t.Run("password too short", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := NewService(mockRepo, jwtManager)

		user, err := svc.Register(ctx, "test@example.com", "short")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.True(t, errors.IsValidationError(err))
	})
}

func TestAuthService_Login(t *testing.T) {
	ctx := context.Background()
	jwtManager := auth.NewJWTManager(config.JWTConfig{
		Secret:     "test-secret",
		Expiration: 24 * 60 * 60,
		Issuer:     "test",
	})

	t.Run("successful login", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := NewService(mockRepo, jwtManager)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		user := &entity.User{
			ID:           uuid.New(),
			Email:        "test@example.com",
			PasswordHash: string(hashedPassword),
		}

		mockRepo.On("FindByEmail", ctx, "test@example.com").Return(user, nil)

		token, err := svc.Login(ctx, "test@example.com", "password123")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := NewService(mockRepo, jwtManager)

		mockRepo.On("FindByEmail", ctx, "notfound@example.com").Return(nil, errors.NewNotFoundError("user"))

		token, err := svc.Login(ctx, "notfound@example.com", "password123")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.True(t, errors.IsUnauthorizedError(err))
		mockRepo.AssertExpectations(t)
	})

	t.Run("wrong password", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := NewService(mockRepo, jwtManager)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
		user := &entity.User{
			ID:           uuid.New(),
			Email:        "test@example.com",
			PasswordHash: string(hashedPassword),
		}

		mockRepo.On("FindByEmail", ctx, "test@example.com").Return(user, nil)

		token, err := svc.Login(ctx, "test@example.com", "wrongpassword")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.True(t, errors.IsUnauthorizedError(err))
		mockRepo.AssertExpectations(t)
	})
}
