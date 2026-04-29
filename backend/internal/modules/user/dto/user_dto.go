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
		ID       string  `json:"id"`
		Name     string  `json:"name"`
		Email    string  `json:"email"`
		RoleName string  `json:"role_name"`
		DetailId *uint   `json:"detail_id"`
		ImageUrl *string `json:"image_url"`
	}

	UserRoleURI struct {
		RoleName string `uri:"role_name" binding:"required,is_non_admin"`
	}
	UserSyncURI struct {
		UserRoleURI
		DetailId uint `uri:"detail_id" binding:"required,gt=0"`
	}

	UserEmailUri struct {
		Email string `uri:"email" binding:"required,email,non_admin_email"`
	}

	UserAdminCreateRequest struct {
		Name     string                `json:"name" form:"name" binding:"required,min=2,max=255"`
		Email    string                `json:"email" form:"email" binding:"required,email"`
		Password string                `json:"password" form:"password" binding:"required,min=8"`
		RoleName string                `json:"role_name" form:"role_kode" binding:"required, custom_role"`
		DetailId *uint                 `json:"detail_id" form:"detail_id" binding:"omitempty,gt=0"`
		Image    *multipart.FileHeader `json:"image" form:"image" binding:"omitempty, custom_ext"`
	}

	UserNonAdminCreateRequest struct {
		Name     string                `json:"name" form:"name" binding:"required,min=2,max=255"`
		Email    string                `json:"email" form:"email" binding:"required,email"`
		Password string                `json:"password" form:"password" binding:"required,min=8"`
		RoleName string                `json:"role_name" form:"role_kode" binding:"required, is_non_admin"`
		DetailId *uint                 `json:"detail_id" form:"detail_id" binding:"required,gt=0"`
		Image    *multipart.FileHeader `json:"image" form:"image" binding:"omitempty, custom_ext"`
	}

	UserAdminUpdateRequest struct {
		Name     string                `json:"name" form:"name" binding:"omitempty,min=2,max=255"`
		Email    string                `json:"email" form:"email" binding:"omitempty,email"`
		Password string                `json:"password" form:"password" binding:"omitempty,min=8"`
		RoleName string                `json:"role_kode" form:"role_kode" binding:"omitempty, custom_role"`
		DetailId *uint                 `json:"detail_id" form:"detail_id" binding:"omitempty,gt=0"`
		Image    *multipart.FileHeader `json:"image" form:"image" binding:"omitempty, custom_ext"`
	}

	UserNonAdminUpdateRequest struct {
		Name     string                `json:"name" form:"name" binding:"omitempty,min=2,max=255"`
		Email    string                `json:"email" form:"email" binding:"omitempty,email"`
		Password string                `json:"password" form:"password" binding:"omitempty,min=8"`
		RoleName string                `json:"role_name" form:"role_name" binding:"omitempty, is_non_admin"`
		DetailId *uint                 `json:"detail_id" form:"detail_id" binding:"omitempty,gt=0"`
		Image    *multipart.FileHeader `json:"image" form:"image" binding:"omitempty, custom_ext"`
	}

	UserLoginRequest struct {
		Email    string `json:"email" form:"email" binding:"required,email"`
		Password string `json:"password" form:"password" binding:"required,min=8"`
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
