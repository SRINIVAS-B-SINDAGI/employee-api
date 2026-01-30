package employee

import (
	"context"
	"testing"

	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/entity"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/valueobject"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/pkg/errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEmployeeRepository struct {
	mock.Mock
}

func (m *MockEmployeeRepository) Create(ctx context.Context, employee *entity.Employee) error {
	args := m.Called(ctx, employee)
	return args.Error(0)
}

func (m *MockEmployeeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Employee, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) Update(ctx context.Context, employee *entity.Employee) error {
	args := m.Called(ctx, employee)
	return args.Error(0)
}

func (m *MockEmployeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEmployeeRepository) GetSalaryStatsByCountry(ctx context.Context, country string) (*valueobject.SalaryStats, error) {
	args := m.Called(ctx, country)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*valueobject.SalaryStats), args.Error(1)
}

func (m *MockEmployeeRepository) GetAvgSalaryByJobTitle(ctx context.Context, jobTitle string) (*valueobject.JobTitleSalaryStats, error) {
	args := m.Called(ctx, jobTitle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*valueobject.JobTitleSalaryStats), args.Error(1)
}

func TestEmployeeService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepository)
		svc := NewService(mockRepo)

		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Employee")).Return(nil)

		emp, err := svc.Create(ctx, "John Doe", "Engineer", "India", decimal.NewFromInt(100000))

		assert.NoError(t, err)
		assert.NotNil(t, emp)
		assert.Equal(t, "John Doe", emp.FullName)
		assert.Equal(t, "Engineer", emp.JobTitle)
		assert.Equal(t, "India", emp.Country)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error - empty name", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepository)
		svc := NewService(mockRepo)

		emp, err := svc.Create(ctx, "", "Engineer", "India", decimal.NewFromInt(100000))

		assert.Error(t, err)
		assert.Nil(t, emp)
		assert.True(t, errors.IsValidationError(err))
	})

	t.Run("validation error - negative salary", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepository)
		svc := NewService(mockRepo)

		emp, err := svc.Create(ctx, "John Doe", "Engineer", "India", decimal.NewFromInt(-100))

		assert.Error(t, err)
		assert.Nil(t, emp)
		assert.True(t, errors.IsValidationError(err))
	})
}

func TestEmployeeService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("found", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepository)
		svc := NewService(mockRepo)

		id := uuid.New()
		expected := &entity.Employee{
			ID:       id,
			FullName: "John Doe",
		}

		mockRepo.On("FindByID", ctx, id).Return(expected, nil)

		emp, err := svc.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expected, emp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepository)
		svc := NewService(mockRepo)

		id := uuid.New()
		mockRepo.On("FindByID", ctx, id).Return(nil, errors.NewNotFoundError("employee"))

		emp, err := svc.GetByID(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, emp)
		assert.True(t, errors.IsNotFoundError(err))
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeService_Update(t *testing.T) {
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepository)
		svc := NewService(mockRepo)

		id := uuid.New()
		existing := &entity.Employee{
			ID:          id,
			FullName:    "John Doe",
			JobTitle:    "Engineer",
			Country:     "India",
			GrossSalary: decimal.NewFromInt(100000),
		}

		mockRepo.On("FindByID", ctx, id).Return(existing, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.Employee")).Return(nil)

		emp, err := svc.Update(ctx, id, "John Updated", "Senior Engineer", "United States", decimal.NewFromInt(150000))

		assert.NoError(t, err)
		assert.NotNil(t, emp)
		assert.Equal(t, "John Updated", emp.FullName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepository)
		svc := NewService(mockRepo)

		id := uuid.New()
		mockRepo.On("FindByID", ctx, id).Return(nil, errors.NewNotFoundError("employee"))

		emp, err := svc.Update(ctx, id, "John", "Engineer", "India", decimal.NewFromInt(100000))

		assert.Error(t, err)
		assert.Nil(t, emp)
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("successful delete", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepository)
		svc := NewService(mockRepo)

		id := uuid.New()
		mockRepo.On("Delete", ctx, id).Return(nil)

		err := svc.Delete(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepository)
		svc := NewService(mockRepo)

		id := uuid.New()
		mockRepo.On("Delete", ctx, id).Return(errors.NewNotFoundError("employee"))

		err := svc.Delete(ctx, id)

		assert.Error(t, err)
		assert.True(t, errors.IsNotFoundError(err))
		mockRepo.AssertExpectations(t)
	})
}
