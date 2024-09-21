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

func OrderStatus(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return (value == "InProgress" || value == "Complete" || value == "Cancel")
}

func RegisterCustomValidations(v *validator.Validate) error {
	if err := v.RegisterValidation("orderStatus", OrderStatus); err != nil {
		return err
	}
	return v.RegisterValidation("minLen", HasMinLength)
}
