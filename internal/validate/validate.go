package validate

import "github.com/go-playground/validator/v10"

// validatorImpl model for validating the bind of requests
type validatorImpl struct {
	v *validator.Validate
}

// Validate method that implements the validator interface for performing validation
func (cv *validatorImpl) Validate(i interface{}) error {
	return cv.v.Struct(i)
}

// Validator validation component interface
type Validator interface {
	Validate(i interface{}) error
}

// New creates a new implementation of the Validator interface
func New() Validator {
	return &validatorImpl{
		v: validator.New(),
	}
}
