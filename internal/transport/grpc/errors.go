package grpc

import (
	"net/http"

	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToGRPCError converts an error to a gRPC status error.
// It maps AppError HTTP status codes to appropriate gRPC status codes.
func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	code := errors.GetStatusCode(err)
	message := errors.GetMessage(err)

	grpcCode := httpStatusToGRPCCode(code)
	return status.Error(grpcCode, message)
}

// httpStatusToGRPCCode maps HTTP status codes to gRPC status codes.
func httpStatusToGRPCCode(httpStatus int) codes.Code {
	switch httpStatus {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusUnprocessableEntity:
		return codes.InvalidArgument
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	default:
		return codes.Unknown
	}
}
