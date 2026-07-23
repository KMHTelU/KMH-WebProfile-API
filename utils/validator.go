package utils

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.Validator.Struct(i)
}

func MapValidationErrors(err error) map[string]string {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok || len(validationErrors) == 0 {
		// Bukan error validasi (mis. JSON malformed) => kembalikan nil agar
		// handler jatuh ke respons "Bad request", bukan "Validation error" kosong.
		return nil
	}
	errorsMap := make(map[string]string)
	for _, e := range validationErrors {
		errorsMap[e.Field()] = e.Tag()
	}
	return errorsMap
}
