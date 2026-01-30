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

func (r *employeeRepository) Update(ctx context.Context, employee *entity.Employee) error {
	result := r.db.WithContext(ctx).Save(employee)
	if result.Error != nil {
		return errors.NewInternalError(result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("employee")
	}
	return nil
}

func (r *employeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entity.Employee{}, "id = ?", id)
	if result.Error != nil {
		return errors.NewInternalError(result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("employee")
	}
	return nil
}
