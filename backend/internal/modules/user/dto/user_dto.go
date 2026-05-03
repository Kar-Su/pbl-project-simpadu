package dto

import (
	"errors"
	"mime/multipart"
	"web-hosting/internal/database/entities"
)

const (
	// FAILED
	MESSAGE_FAILED_REGISTER_USER   = "failed to register user"
	MESSAGE_FAILED_LOGIN_USER      = "failed to login user"
	MESSAGE_FAILED_GET_USER        = "failed to get user"
	MESSAGE_FAILED_GET_LIST_USER   = "failed to get list of users"
	MESSAGE_FAILED_UPDATE_USER     = "failed to update user"
	MESSAGE_FAILED_DELETE_USER     = "failed to delete user"
	MESSAGE_FAILED_PROSES_REQUEST  = "failed to process request"
	MESSAGE_FAILED_TOKEN_NOT_FOUND = "token not found"
	MESSAGE_FAILED_TOKEN_INVALID   = "token invalid"
	MESSAGE_FAILED_DENIED_ACCESS   = "access denied"
	MESSAGE_FAILED_BAD_REQUEST     = "bad request"

	// SUCCESS
	MESSAGE_SUCCESS_REGISTER_USER  = "user registered successfully"
	MESSAGE_SUCCESS_LOGIN_USER     = "user logged in successfully"
	MESSAGE_SUCCESS_GET_USER       = "user retrieved successfully"
	MESSAGE_SUCCESS_GET_LIST_USER  = "list of users retrieved successfully"
	MESSAGE_SUCCESS_UPDATE_USER    = "user updated successfully"
	MESSAGE_SUCCESS_DELETE_USER    = "user deleted successfully"
	MESSAGE_SUCCESS_PROSES_REQUEST = "request processed successfully"
)

var (
	ErrCreateUser = errors.New("failed to create user")

	ErrUpdateUser         = errors.New("failed to update user")
	ErrDeleteUser         = errors.New("failed to delete user")
	ErrLoginUser          = errors.New("failed to login user")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidAdminRole   = errors.New("invalid not admin role")
	ErrRoleNotFound       = errors.New("role not found")
	ErrTokenInvalid       = errors.New("token invalid")
	ErrTokenExpired       = errors.New("token expired")
)

type (
	UserResponse struct {
		ID       string  `json:"id" example:"ijq0rq0jfa //(UUID v7)"`
		Name     string  `json:"name" example:"rezi"`
		Email    string  `json:"email" example:"rezi@example.com"`
		RoleName string  `json:"role_name" example:"raja-sawit"`
		DetailId *uint   `json:"detail_id" example:"1"`
		ImageUrl *string `json:"image_url" example:"path/to/image.jpg"`
	}

	UserRoleURI struct {
		RoleName string `uri:"role_name" binding:"required,is_non_admin"`
	}
	UserSyncURI struct {
		UserRoleURI
		DetailId uint `uri:"detail_id" binding:"required,gt=0"`
	}

	UserEmailRequest struct {
		Email string `form:"email" binding:"required,email,non_admin_email"`
	}

	UserAdminCreateRequest struct {
		Name     string                `json:"name" form:"name" binding:"required,min=2,max=255" example:"rezi // required, min 2 characters, max 255 characters"`
		Email    string                `json:"email" form:"email" binding:"required,email" example:"rezi@example.com // required, must be a valid email address"`
		Password string                `json:"password" form:"password" binding:"required,min=8" example:"inipasswordrezi // required, min 8 characters"`
		RoleName string                `json:"role_name" form:"role_kode" binding:"required" example:"raja-nyawit // required, must be a valid role name"`
		DetailId *uint                 `json:"detail_id" form:"detail_id" binding:"omitempty,gt=0" example:"1"`
		Image    *multipart.FileHeader `json:"image" form:"image" binding:"omitempty,custom_ext" example:"path/to/image.jpg"`
	}

	UserNonAdminCreateRequest struct {
		Name     string                `json:"name" form:"name" binding:"required,min=2,max=255" example:"Rezi // required, min 2 max 255 characters"`
		Email    string                `json:"email" form:"email" binding:"required,email" example:"rezi@example.com // required, must be a valid email address"`
		Password string                `json:"password" form:"password" binding:"required,min=8" example:"inipasswordrezi // required, min 8 characters"`
		RoleName string                `json:"role_name" form:"role_kode" binding:"required,is_non_admin" example:"raja-nyawit // required, must be a valid role name"`
		DetailId *uint                 `json:"detail_id" form:"detail_id" binding:"required,gt=0" example:"1"`
		Image    *multipart.FileHeader `json:"image" form:"image" binding:"omitempty,custom_ext" example:"path/to/image.jpg"`
	}

	UserAdminUpdateRequest struct {
		Name     string                `json:"name" form:"name" binding:"omitempty,min=2,max=255" example:"Rezi // optional, min 2 max 255 characters"`
		Email    string                `json:"email" form:"email" binding:"omitempty,email" example:"rezi@example.com // optional, must be a valid email address"`
		Password string                `json:"password" form:"password" binding:"omitempty,min=8" example:"inipasswordrezi // optional, min 8 characters"`
		RoleName string                `json:"role_name" form:"role_name" binding:"omitempty" example:"raja-nyawit // optional"`
		DetailId *uint                 `json:"detail_id" form:"detail_id" binding:"omitempty,gt=0" example:"1"`
		Image    *multipart.FileHeader `json:"image" form:"image" binding:"omitempty,custom_ext" example:"path/to/image.jpg"`
	}

	UserNonAdminUpdateRequest struct {
		Name     string                `json:"name" form:"name" binding:"omitempty,min=2,max=255" example:"rezi"`
		Email    string                `json:"email" form:"email" binding:"omitempty,email" example:"rezi@example.com // optional, must be a valid email address"`
		Password string                `json:"password" form:"password" binding:"omitempty,min=8" example:"inipasswordrezi // optional, min 8 characters"`
		RoleName string                `json:"role_name" form:"role_name" binding:"omitempty,is_non_admin" example:"raja-nyawit // optional"`
		DetailId *uint                 `json:"detail_id" form:"detail_id" binding:"omitempty,gt=0" example:"1"`
		Image    *multipart.FileHeader `json:"image" form:"image" binding:"omitempty,custom_ext" example:"path/to/image.jpg"`
	}

	UserLoginRequest struct {
		Email    string `json:"email" form:"email" binding:"required,email" example:"rezi@example.com // required, must be a valid email address"`
		Password string `json:"password" form:"password" binding:"required,min=8" example:"inipasswordrezi // required, min 8 characters"`
	}
)

func ToUserResponse(user entities.User) UserResponse {
	return UserResponse{
		ID:       user.ID.String(),
		Name:     user.Name,
		Email:    user.Email,
		RoleName: user.Role.Name,
		DetailId: user.DetailID,
		ImageUrl: user.ImageUrl,
	}
}
