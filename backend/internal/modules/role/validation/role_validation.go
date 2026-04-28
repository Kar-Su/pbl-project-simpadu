package validation

import (
	"slices"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/helpers"

	"github.com/go-playground/validator/v10"
)

type UserValidation struct {
	validate *validator.Validate
}

func NewRoleValidation() *UserValidation {
	validate := validator.New()
	validate.SetTagName("binding")
	validate.RegisterValidation("role_exist", existRoleName)

	return &UserValidation{
		validate: validate,
	}
}

func existRoleName(fl validator.FieldLevel) bool {
	roleName := helpers.NormalizeString(fl.Field().String())

	existRoleName := []string{constants.ROLE_ADMIN,
		constants.ROLE_DOSEN,
		constants.ROLE_MAHASISWA,
		constants.ROLE_SUPER_ADMIN,
		constants.ROLE_ADMIN,
	}

	isExist := slices.Contains(existRoleName, roleName)
	return isExist
}
