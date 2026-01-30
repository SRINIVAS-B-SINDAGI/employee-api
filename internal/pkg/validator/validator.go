package validator

import (
	"regexp"
	"strings"

	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/pkg/errors"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.NewValidationError("email is required")
	}
	if !emailRegex.MatchString(email) {
		return errors.NewValidationError("invalid email format")
	}
	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.NewValidationError("password is required")
	}
	if len(password) < 8 {
		return errors.NewValidationError("password must be at least 8 characters")
	}
	return nil
}

func ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return errors.NewValidationError(fieldName + " is required")
	}
	return nil
}
