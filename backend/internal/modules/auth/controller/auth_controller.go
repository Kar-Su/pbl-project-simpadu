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
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REFRESH_TOKEN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REFRESH_TOKEN, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *authController) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordRequest

	if err := ctx.ShouldBindJSON(req); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
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
