package validators

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func HasMinLength(fl validator.FieldLevel) bool {
	param := fl.Param()
	minLen, err := strconv.Atoi(param)
	if err != nil {
		return false
	}

	value, ok := fl.Field().Interface().([]uuid.UUID)
	return ok && len(value) >= minLen
}

func RegisterCustomValidations(v *validator.Validate) error {
	return v.RegisterValidation("minLen", HasMinLength)
}
