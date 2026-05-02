package dto

import (
	"errors"
)

const (
	MESSAGE_FAILED_REFRESH_TOKEN       = "failed refresh token"
	MESSAGE_FAILED_LOGOUT              = "failed logout"
	MESSAGE_FAILED_SEND_PASSWORD_RESET = "failed send password reset"
	MESSAGE_FAILED_RESET_PASSWORD      = "failed reset password"
	MESSAGE_FAILED_RETRIEVE_SECRET_KEY = "failed retrieve secret key"
	MESSAGE_FAILED_FIND_REFRESH_TOKEN  = "failed find refresh token"

	MESSAGE_SUCCESS_REFRESH_TOKEN       = "success refresh token"
	MESSAGE_SUCCESS_LOGOUT              = "success logout"
	MESSAGE_SUCCESS_SEND_PASSWORD_RESET = "success send password reset"
	MESSAGE_SUCCESS_RESET_PASSWORD      = "success reset password"
	MESSAGE_SUCCESS_FIND_REFRESH_TOKEN  = "success find refresh token"
)

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrInvalidCredentials   = errors.New("invalid credentials")
)

type (
	RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	TokenResponse struct {
		UserName     string `json:"user_name"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		RoleName     string `json:"role_name"`
	}

	ResetPasswordRequest struct {
		Email       string `json:"email" binding:"required,email"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	RefreshTokenResponse struct {
		RefreshToken string `json:"refresh_token"`
		Exp          int64  `json:"expired_at"`
	}
)
