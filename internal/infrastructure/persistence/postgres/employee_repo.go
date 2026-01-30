package postgres

import (
	"context"

	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/entity"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/repository"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) repository.EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) Create(ctx context.Context, employee *entity.Employee) error {
	if err := r.db.WithContext(ctx).Create(employee).Error; err != nil {
		return errors.NewInternalError(err)
	}
	return nil
}

func (r *employeeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Employee, error) {
	var employee entity.Employee
	if err := r.db.WithContext(ctx).First(&employee, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("employee")
		}
		return nil, errors.NewInternalError(err)
	}
	return &employee, nil
}
