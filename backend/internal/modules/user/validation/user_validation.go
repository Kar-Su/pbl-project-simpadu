package validation

import (
	"mime/multipart"
	"path/filepath"
	"slices"
	"strings"
	"web-hosting/internal/modules/user/dto"

	"github.com/go-playground/validator/v10"
)

type UserValidation struct {
	validate *validator.Validate
}

func NewUserValidation() *UserValidation {
	validate := validator.New()
	validate.SetTagName("binding")
	validate.RegisterValidation("is_non_admin", validateNonAdmin)
	validate.RegisterValidation("custom_ext", customExtImage)

	return &UserValidation{
		validate: validate,
	}
}

func (v *UserValidation) ValidateUpdateAdminRequest(req dto.UserAdminUpdateRequest) error {
	return v.validate.Struct(req)
}

func (v *UserValidation) ValidateUpdateNonAdminRequest(req dto.UserNonAdminUpdateRequest) error {
	return v.validate.Struct(req)
}

func validateNonAdmin(fl validator.FieldLevel) bool {
	roleId := fl.Field().String()

	return roleId != "r-sa" && roleId != "r-ad"
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
