package validation

import (
	"mime/multipart"
	"path/filepath"
	"slices"
	"strings"
	"web-hosting/internal/modules/user/dto"
	"web-hosting/internal/package/constants"
	help "web-hosting/internal/package/helpers"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type UserValidation struct {
	validate *validator.Validate
}

func NewUserValidation() *UserValidation {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.SetTagName("binding")
		v.RegisterValidation("is_non_admin", validateRoleNonAdmin)
		v.RegisterValidation("custom_ext", customExtImage)
		v.RegisterValidation("non_admin_email", validateNonAdminEmail)
	}
	return &UserValidation{
		validate: v,
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

func validateRoleNonAdmin(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	role = help.NormalizeString(role)

	return role != constants.ROLE_ADMIN_PEGAWAI && role != constants.ROLE_ADMIN_MAHASISWA && role != constants.ROLE_ADMIN_KEUANGAN && role != constants.ROLE_SUPER_ADMIN && role != constants.ROLE_DOSEN && role != constants.ROLE_ADMIN_AKADEMIK
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

func validateNonAdminEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	return email != constants.EMAIL_ADMIN_PEGAWAI && email != constants.EMAIL_ADMIN_MAHASISWA && email != constants.EMAIL_ADMIN_KEUANGAN && email != constants.EMAIL_SUPER_ADMIN
}
