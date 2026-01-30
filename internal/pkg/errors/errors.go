package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: message,
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "internal server error",
		Err:     err,
	}
}

func GetStatusCode(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return http.StatusInternalServerError
}

func GetMessage(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Message
	}
	return "internal server error"
}

func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

func IsNotFoundError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == http.StatusNotFound
	}
	return false
}

func NewUnauthorizedError(message string) *AppError {
	if message == "" {
		message = "unauthorized"
	}
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}
