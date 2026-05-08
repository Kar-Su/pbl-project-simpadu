package validation

import (
	"slices"
	"web-hosting/internal/package/helpers"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type TahunAkademikValidator struct {
	validate *validator.Validate
}

func NewTahunAkademikValidator() *TahunAkademikValidator {
	v, ok := binding.Validator.Engine().(*validator.Validate)

	if ok {
		v.SetTagName("binding")
		v.RegisterValidation("enumTipeSemester", validateEnumTipeSemester)
		v.RegisterValidation("enumStatus", validateEnumStatus)
	}

	return &TahunAkademikValidator{
		validate: v,
	}
}

func validateEnumTipeSemester(field validator.FieldLevel) bool {
	req := helpers.NormalizeString(field.Field().String())

	enum := []string{"ganjil", "genap"}

	if slices.Contains(enum, req) {
		return true
	}

	return false
}

func validateEnumStatus(field validator.FieldLevel) bool {
	req := helpers.NormalizeString(field.Field().String())

	enum := []string{"aktif", "nonaktif"}

	if slices.Contains(enum, req) {
		return true
	}

	return false
}
