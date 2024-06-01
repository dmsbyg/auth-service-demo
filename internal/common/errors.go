package common

import (
	"errors"
	"fmt"
)

var (
	ErrUnauthorized       error = errors.New("unauthorized")
	ErrInvalidCredentials error = errors.New("email or password are incorrect")
)

// InternalAppError is an erorr that should not be returned to API's user
type InternalAppError struct {
	Err error
}

func (i InternalAppError) Error() string {
	return "internal app error"
}

func NewInternalAppError(err error) error {
	return InternalAppError{Err: err}
}

type DuplicateError struct {
	Entity string
}

func (e DuplicateError) Error() string {
	return fmt.Sprintf("duplicate entry for %s", e.Entity)
}
