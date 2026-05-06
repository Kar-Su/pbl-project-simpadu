package validation

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type ProdiValidation struct {
	validate *validator.Validate
}

func NewProdiValidation() *ProdiValidation {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.SetTagName("binding")
		v.RegisterValidation("enum_jenjang", validateJenjang)
	}

	return &ProdiValidation{
		validate: v,
	}
}

func validateJenjang(fl validator.FieldLevel) bool {
	jenjang := fl.Field().String()

	return jenjang == "D3" || jenjang == "D4"
}
