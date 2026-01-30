package postgres

import (
	"context"

	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/entity"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/repository"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/valueobject"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/pkg/errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

func (r *employeeRepository) GetSalaryStatsByCountry(ctx context.Context, country string) (*valueobject.SalaryStats, error) {
	var result struct {
		MinSalary decimal.Decimal
		MaxSalary decimal.Decimal
		AvgSalary decimal.Decimal
		Count     int64
	}

	err := r.db.WithContext(ctx).
		Model(&entity.Employee{}).
		Select("MIN(gross_salary) as min_salary, MAX(gross_salary) as max_salary, AVG(gross_salary) as avg_salary, COUNT(*) as count").
		Where("country = ?", country).
		Scan(&result).Error

	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	if result.Count == 0 {
		return nil, errors.NewNotFoundError("employees in country")
	}

	return &valueobject.SalaryStats{
		MinSalary: result.MinSalary,
		MaxSalary: result.MaxSalary,
		AvgSalary: result.AvgSalary.Round(2),
		Count:     result.Count,
	}, nil
}

func (r *employeeRepository) GetAvgSalaryByJobTitle(ctx context.Context, jobTitle string) (*valueobject.JobTitleSalaryStats, error) {
	var result struct {
		AvgSalary decimal.Decimal
		Count     int64
	}

	err := r.db.WithContext(ctx).
		Model(&entity.Employee{}).
		Select("AVG(gross_salary) as avg_salary, COUNT(*) as count").
		Where("job_title = ?", jobTitle).
		Scan(&result).Error

	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	if result.Count == 0 {
		return nil, errors.NewNotFoundError("employees with job title")
	}

	return &valueobject.JobTitleSalaryStats{
		JobTitle:  jobTitle,
		AvgSalary: result.AvgSalary.Round(2),
		Count:     result.Count,
	}, nil
}
