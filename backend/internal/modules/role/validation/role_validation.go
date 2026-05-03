package validation

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type UserValidation struct {
	validate *validator.Validate
}

func NewRoleValidation() *UserValidation {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.SetTagName("binding")
	}

	return &UserValidation{
		validate: v,
	}
}
