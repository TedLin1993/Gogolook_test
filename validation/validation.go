package validation

import (
	. "SimpleApi/model"

	"github.com/go-playground/validator/v10"
)

func Init() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("validStatus", func(fl validator.FieldLevel) bool {
		status, ok := fl.Field().Interface().(Status)
		if !ok {
			return false
		}
		switch status {
		case Incomplete, Completed:
			return true
		default:
			return false
		}
	})
	return validate
}
