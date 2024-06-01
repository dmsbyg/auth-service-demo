package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func New() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) Validate(i interface{}) error {
	err := v.validator.Struct(i)
	if errs, ok := err.(validator.ValidationErrors); ok { //nolint:errorlint
		for _, err := range errs {
			switch err.Tag() {
			case "required":
				return fmt.Errorf("'%s' field is required", err.Field())
			case "email":
				return fmt.Errorf("'%s' field shall be a valid email", err.Field())
			case "gte":
				return fmt.Errorf("'%s' field shall be equal to or greater than %s", err.Field(), err.Param())
			case "gt":
				return fmt.Errorf("'%s' field shall be greater than %s", err.Field(), err.Param())
			case "lte":
				return fmt.Errorf("'%s' field shall be equal to or lower than %s", err.Field(), err.Param())
			case "lt":
				return fmt.Errorf("'%s' field shall be lower than %s", err.Field(), err.Param())
			case "min":
				return fmt.Errorf("'%s' field shall be minimum of %s characters", err.Field(), err.Param())
			case "max":
				return fmt.Errorf("'%s' field shall be maximum of %s characters", err.Field(), err.Param())
			}
		}
	}
	return err
}
