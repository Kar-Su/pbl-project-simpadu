package validation

import (
	"slices"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type PresensiValidation struct {
	validate *validator.Validate
}

func NewPresensiValidation() *PresensiValidation {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.SetTagName("binding")
		v.RegisterValidation("enum_status", validateEnumStatus)
	}

	return &PresensiValidation{
		validate: v,
	}
}

func validateEnumStatus(fl validator.FieldLevel) bool {
	status := strings.ToLower(fl.Field().String())

	validStatus := []string{"hadir", "izin", "sakit", "alpha"}
	if slices.Contains(validStatus, status) {
		return true
	}

	return false
}
