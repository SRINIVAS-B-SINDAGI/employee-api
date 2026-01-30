package employee

import (
	"context"

	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/entity"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/repository"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/pkg/errors"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/pkg/validator"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Service interface {
	Create(ctx context.Context, fullName, jobTitle, country string, grossSalary decimal.Decimal) (*entity.Employee, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Employee, error)
	Update(ctx context.Context, id uuid.UUID, fullName, jobTitle, country string, grossSalary decimal.Decimal) (*entity.Employee, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo repository.EmployeeRepository
}

func NewService(repo repository.EmployeeRepository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, fullName, jobTitle, country string, grossSalary decimal.Decimal) (*entity.Employee, error) {
	if err := s.validateEmployee(fullName, jobTitle, country, grossSalary); err != nil {
		return nil, err
	}

	employee := entity.NewEmployee(fullName, jobTitle, country, grossSalary)
	if err := s.repo.Create(ctx, employee); err != nil {
		return nil, err
	}

	return employee, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*entity.Employee, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) validateEmployee(fullName, jobTitle, country string, grossSalary decimal.Decimal) error {
	if err := validator.ValidateRequired(fullName, "full_name"); err != nil {
		return err
	}
	if err := validator.ValidateRequired(jobTitle, "job_title"); err != nil {
		return err
	}
	if err := validator.ValidateRequired(country, "country"); err != nil {
		return err
	}
	if grossSalary.LessThan(decimal.Zero) {
		return errors.NewValidationError("gross_salary cannot be negative")
	}
	return nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, fullName, jobTitle, country string, grossSalary decimal.Decimal) (*entity.Employee, error) {
	if err := s.validateEmployee(fullName, jobTitle, country, grossSalary); err != nil {
		return nil, err
	}

	employee, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	employee.FullName = fullName
	employee.JobTitle = jobTitle
	employee.Country = country
	employee.GrossSalary = grossSalary

	if err := s.repo.Update(ctx, employee); err != nil {
		return nil, err
	}

	return employee, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
