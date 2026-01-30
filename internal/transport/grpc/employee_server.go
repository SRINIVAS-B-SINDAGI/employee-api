package grpc

import (
	"context"

	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/domain/entity"
	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/pkg/errors"
	employeeuc "github.com/SRINIVAS-B-SINDAGI/employee-api/internal/usecase/employee"
	employeev1 "github.com/SRINIVAS-B-SINDAGI/employee-api/proto/employee/v1"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type employeeServer struct {
	employeev1.UnimplementedEmployeeServiceServer
	service employeeuc.Service
}

func NewEmployeeServer(service employeeuc.Service) employeev1.EmployeeServiceServer {
	return &employeeServer{
		service: service,
	}
}

func (s *employeeServer) CreateEmployee(ctx context.Context, req *employeev1.CreateEmployeeRequest) (*employeev1.Employee, error) {
	grossSalary, err := decimal.NewFromString(req.GetGrossSalary())
	if err != nil {
		return nil, ToGRPCError(errors.NewValidationError("invalid gross_salary format"))
	}

	employee, err := s.service.Create(ctx, req.GetFullName(), req.GetJobTitle(), req.GetCountry(), grossSalary)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return entityToProto(employee), nil
}

func (s *employeeServer) GetEmployee(ctx context.Context, req *employeev1.GetEmployeeRequest) (*employeev1.Employee, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, ToGRPCError(errors.NewValidationError("invalid employee id format"))
	}

	employee, err := s.service.GetByID(ctx, id)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return entityToProto(employee), nil
}

func entityToProto(e *entity.Employee) *employeev1.Employee {
	return &employeev1.Employee{
		Id:          e.ID.String(),
		FullName:    e.FullName,
		JobTitle:    e.JobTitle,
		Country:     e.Country,
		GrossSalary: e.GrossSalary.String(),
		CreatedAt:   timestamppb.New(e.CreatedAt),
		UpdatedAt:   timestamppb.New(e.UpdatedAt),
	}
}
