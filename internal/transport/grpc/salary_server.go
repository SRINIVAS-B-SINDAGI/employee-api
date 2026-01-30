package grpc

import (
	"context"

	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/pkg/errors"
	salaryuc "github.com/SRINIVAS-B-SINDAGI/employee-api/internal/usecase/salary"
	salaryv1 "github.com/SRINIVAS-B-SINDAGI/employee-api/proto/salary/v1"
	"github.com/google/uuid"
)

// salaryServer implements the SalaryServiceServer interface.
type salaryServer struct {
	salaryv1.UnimplementedSalaryServiceServer
	service salaryuc.Service
}

// NewSalaryServer creates a new gRPC salary server.
func NewSalaryServer(service salaryuc.Service) salaryv1.SalaryServiceServer {
	return &salaryServer{
		service: service,
	}
}

// CalculateNetSalary calculates the net salary for an employee.
func (s *salaryServer) CalculateNetSalary(ctx context.Context, req *salaryv1.CalculateNetSalaryRequest) (*salaryv1.CalculateNetSalaryResponse, error) {
	employeeID, err := uuid.Parse(req.GetEmployeeId())
	if err != nil {
		return nil, ToGRPCError(errors.NewValidationError("invalid employee_id format"))
	}

	salary, err := s.service.CalculateNetSalary(ctx, employeeID)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &salaryv1.CalculateNetSalaryResponse{
		GrossSalary: salary.GrossSalary.String(),
		TaxRate:     salary.TaxRate.String(),
		TaxAmount:   salary.TaxAmount.String(),
		NetSalary:   salary.NetSalary.String(),
	}, nil
}

// GetSalaryStatsByCountry returns salary statistics for a country.
func (s *salaryServer) GetSalaryStatsByCountry(ctx context.Context, req *salaryv1.GetSalaryStatsByCountryRequest) (*salaryv1.SalaryStatsResponse, error) {
	stats, err := s.service.GetSalaryStatsByCountry(ctx, req.GetCountry())
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &salaryv1.SalaryStatsResponse{
		MinSalary: stats.MinSalary.String(),
		MaxSalary: stats.MaxSalary.String(),
		AvgSalary: stats.AvgSalary.String(),
		Count:     stats.Count,
	}, nil
}

// GetAvgSalaryByJobTitle returns average salary for a job title.
func (s *salaryServer) GetAvgSalaryByJobTitle(ctx context.Context, req *salaryv1.GetAvgSalaryByJobTitleRequest) (*salaryv1.JobTitleSalaryStatsResponse, error) {
	stats, err := s.service.GetAvgSalaryByJobTitle(ctx, req.GetJobTitle())
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &salaryv1.JobTitleSalaryStatsResponse{
		JobTitle:  stats.JobTitle,
		AvgSalary: stats.AvgSalary.String(),
		Count:     stats.Count,
	}, nil
}
