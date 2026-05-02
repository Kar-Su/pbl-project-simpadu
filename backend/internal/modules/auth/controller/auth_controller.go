package controller

import (
	"errors"
	"net/http"
	"web-hosting/internal/modules/auth/dto"
	authServ "web-hosting/internal/modules/auth/service"
	userDto "web-hosting/internal/modules/user/dto"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type AuthController interface {
	FindRefreshToken(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
}

type authController struct {
	authService authServ.AuthService
	db          *gorm.DB
}

func NewAuthController(injector do.Injector, authService authServ.AuthService, db *gorm.DB) AuthController {
	return &authController{
		authService: authService,
		db:          db,
	}
}

// FindRefreshToken godoc
// @Summary      Cari Detail Refresh Token
// @Description  Mengambil data detail dari sebuah refresh token berdasarkan string token
// @Description Access:
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refresh_token  path      string  true  "Refresh Token"
// @Success      200  {object}  utils.Response{data=dto.RefreshTokenResponse}
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /api/auth/refresh-token/{refresh_token} [get]
func (c *authController) FindRefreshToken(ctx *gin.Context) {
	token := ctx.Param("refresh_token")

	result, err := c.authService.FindRefreshToken(ctx.Request.Context(), token)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_FIND_REFRESH_TOKEN, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_FIND_REFRESH_TOKEN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_FIND_REFRESH_TOKEN, result)
	ctx.JSON(http.StatusOK, res)
}

// Login godoc
// @Summary      User Login
// @Description  Proses autentikasi user untuk mendapatkan Access Token dan Refresh Token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body userDto.UserLoginRequest  true  "Payload Login"
// @Success      200  {object}  utils.Response{data=dto.TokenResponse}
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /api/auth/login [post]
func (c *authController) Login(ctx *gin.Context) {
	var req userDto.UserLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.authService.Login(ctx, req)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(userDto.MESSAGE_FAILED_LOGIN_USER, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(userDto.MESSAGE_FAILED_LOGIN_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(userDto.MESSAGE_SUCCESS_LOGIN_USER, result)
	ctx.JSON(http.StatusOK, res)
}

// Logout godoc
// @Summary      User Logout
// @Description  Menghapus session user dan menonaktifkan token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /api/auth/logout [post]
func (c *authController) Logout(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)

	if err := c.authService.Logout(ctx, userId); err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_LOGOUT, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_LOGOUT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_LOGOUT, nil)
	ctx.JSON(http.StatusOK, res)
}

// RefreshToken godoc
// @Summary      Refresh Access Token
// @Description  Mendapatkan access token baru menggunakan refresh token yang masih valid
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RefreshTokenRequest  true  "Payload Refresh Token"
// @Success      200  {object}  utils.Response{data=dto.TokenResponse}
// @Failure      401  {object}  utils.Response
// @Failure      400  {object}  utils.ResponseErr
// @Failure      500  {object}  utils.ResponseErr
// @Router       /api/auth/refresh-token [post]
func (c *authController) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.authService.RefreshToken(ctx, req)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REFRESH_TOKEN, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		if errors.Is(err, dto.ErrRefreshTokenExpired) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REFRESH_TOKEN, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REFRESH_TOKEN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REFRESH_TOKEN, result)
	ctx.JSON(http.StatusOK, res)
}

// ResetPassword godoc
// @Summary      Reset Password User
// @Description  Mengubah password user. Hanya bisa dilakukan pemilik akun atau Super Admin.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body 	dto.ResetPasswordRequest  true  "Payload Reset Password"
// @Success      200  {object}  utils.Response
// @Failure      401  {object}  utils.ResponseErr
// @Failure      400  {object}  utils.ResponseErr
// @Failure      500  {object}  utils.ResponseErr
// @Router       /api/auth/reset-password [post]
func (c *authController) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userEmail := ctx.MustGet("user_email").(string)
	userRole := ctx.MustGet("role_name").(string)
	if userEmail != req.Email && userRole != constants.ROLE_SUPER_ADMIN {
		res := utils.BuildResponseFailed("User unauthorized", "You are not authorized to reset this password", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}

	if err := c.authService.ResetPassword(ctx, req); err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_SEND_PASSWORD_RESET, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_SEND_PASSWORD_RESET, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_RESET_PASSWORD, nil)
	ctx.JSON(http.StatusOK, res)
}
