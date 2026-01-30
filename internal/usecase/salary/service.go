package salary

import (
	"context"

	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/repository"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/valueobject"
	"github.com/google/uuid"
)

type Service interface {
	CalculateNetSalary(ctx context.Context, employeeID uuid.UUID) (*valueobject.Salary, error)
	GetSalaryStatsByCountry(ctx context.Context, country string) (*valueobject.SalaryStats, error)
	GetAvgSalaryByJobTitle(ctx context.Context, jobTitle string) (*valueobject.JobTitleSalaryStats, error)
}

type service struct {
	employeeRepo repository.EmployeeRepository
}

func NewService(employeeRepo repository.EmployeeRepository) Service {
	return &service{employeeRepo: employeeRepo}
}

func (s *service) CalculateNetSalary(ctx context.Context, employeeID uuid.UUID) (*valueobject.Salary, error) {
	employee, err := s.employeeRepo.FindByID(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	salary := valueobject.CalculateNetSalary(
		employee.GrossSalary,
		valueobject.Country(employee.Country),
	)

	return &salary, nil
}

func (s *service) GetSalaryStatsByCountry(ctx context.Context, country string) (*valueobject.SalaryStats, error) {
	return s.employeeRepo.GetSalaryStatsByCountry(ctx, country)
}

func (s *service) GetAvgSalaryByJobTitle(ctx context.Context, jobTitle string) (*valueobject.JobTitleSalaryStats, error) {
	return s.employeeRepo.GetAvgSalaryByJobTitle(ctx, jobTitle)
}
