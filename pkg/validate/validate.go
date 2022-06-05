package validate

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

func ValidateVar(tag string, fields ...interface{}) bool {
	for _, obj := range fields {
		if err := validate.Var(obj, tag); err != nil {
			return false
		}
	}
	return true
}

func ValidateGtInt(val int, fields ...interface{}) bool {
	tag := fmt.Sprintf("required,gt=%d", val)
	return ValidateVar(tag, fields...)
}

func ValidateStringEmpty(fields ...interface{}) bool {
	return ValidateVar("required", fields...)
}

func ValidateEmail(fields ...interface{}) bool {
	return ValidateVar("required,email", fields...)
}
