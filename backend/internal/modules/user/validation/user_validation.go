package validation

import (
	"mime/multipart"
	"path/filepath"
	"slices"
	"strings"
	"web-hosting/internal/modules/user/dto"
	"web-hosting/internal/package/constants"
	help "web-hosting/internal/package/helpers"

	"github.com/go-playground/validator/v10"
)

type UserValidation struct {
	validate *validator.Validate
}

func NewUserValidation() *UserValidation {
	validate := validator.New()
	validate.SetTagName("binding")
	validate.RegisterValidation("custom_role", validateRole)
	validate.RegisterValidation("is_non_admin", validateRoleNonAdmin)
	validate.RegisterValidation("custom_ext", customExtImage)

	return &UserValidation{
		validate: validate,
	}
}

func (v *UserValidation) ValidateSyncURI(req dto.UserSyncURI) error {
	return v.validate.Struct(req)
}

func (v *UserValidation) ValidateUserRoleURI(req dto.UserRoleURI) error {
	return v.validate.Struct(req)
}

func (v *UserValidation) ValidateRegisterRequest(req dto.UserAdminCreateRequest) error {
	return v.validate.Struct(req)
}

func (v *UserValidation) ValidateUpdateAdminRequest(req dto.UserAdminUpdateRequest) error {
	return v.validate.Struct(req)
}

func (v *UserValidation) ValidateUpdateNonAdminRequest(req dto.UserNonAdminUpdateRequest) error {
	return v.validate.Struct(req)
}

func validateRole(fl validator.FieldLevel) bool {
	role := fl.Field().String()

	role = help.NormalizeString(role)

	return role == constants.ROLE_ADMIN || role == constants.ROLE_SUPER_ADMIN || role == constants.ROLE_MAHASISWA || role == constants.ROLE_DOSEN
}

func validateRoleNonAdmin(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	role = help.NormalizeString(role)

	return role != constants.ROLE_ADMIN && role != constants.ROLE_SUPER_ADMIN
}

func customExtImage(fl validator.FieldLevel) bool {
	image, ok := fl.Field().Interface().(*multipart.FileHeader)
	if !ok || image == nil {
		return true
	}
	if image.Filename == "null" {
		return true
	}

	ext := strings.ToLower(filepath.Ext(image.Filename))
	validExts := []string{".jpg", ".jpeg", ".png"}

	return slices.Contains(validExts, ext)
}
