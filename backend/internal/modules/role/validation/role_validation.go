package validation

import (
	"github.com/go-playground/validator/v10"
)

type UserValidation struct {
	validate *validator.Validate
}

func NewRoleValidation() *UserValidation {
	validate := validator.New()
	validate.SetTagName("binding")

	return &UserValidation{
		validate: validate,
	}
}
