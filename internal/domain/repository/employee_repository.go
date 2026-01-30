package repository

import (
	"context"

	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/entity"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/valueobject"
	"github.com/google/uuid"
)

type EmployeeRepository interface {
	Create(ctx context.Context, employee *entity.Employee) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Employee, error)
	Update(ctx context.Context, employee *entity.Employee) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetSalaryStatsByCountry(ctx context.Context, country string) (*valueobject.SalaryStats, error)
	GetAvgSalaryByJobTitle(ctx context.Context, jobTitle string) (*valueobject.JobTitleSalaryStats, error)
}
