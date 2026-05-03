package controller

import (
	"errors"
	"net/http"
	"web-hosting/internal/modules/auth/dto"
	authServ "web-hosting/internal/modules/auth/service"
	userDto "web-hosting/internal/modules/user/dto"
	"web-hosting/internal/package/constants"
	_ "web-hosting/internal/package/swagger"
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
// @Description  Mengambil data detail dari sebuah refresh token berdasarkan string token-nya.
// @Description
// @Description  **Akses:** Authenticated User (memerlukan Bearer Token di header Authorization).
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `404` Token tidak ditemukan di database -> `message: "failed find refresh token", error: "refresh token not found"`
// @Description  - `500` Kesalahan internal server -> `message: "failed find refresh token", error: "Internal Error"`
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        refresh_token  path      string  true  "String Refresh Token"  example(MBG-JAYA67)
// @Success      200  {object}  utils.Response[dto.RefreshTokenResponse,any]
// @Failure      401  {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      404  {object}  swagger.ErrFindRefreshTokenNotFound
// @Failure      500  {object}  swagger.ErrFindRefreshTokenInternal
// @Router       /api/auth/refresh-token/{refresh_token} [get]
func (c *authController) FindRefreshToken(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	token := ctx.Param("refresh_token")

	result, err := c.authService.FindRefreshToken(ctx.Request.Context(), token)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_FIND_REFRESH_TOKEN, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_FIND_REFRESH_TOKEN, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_FIND_REFRESH_TOKEN, result, path)
	ctx.JSON(http.StatusOK, res)
}

// Login godoc
// @Summary      User Login
// @Description  Proses autentikasi user untuk mendapatkan Access Token dan Refresh Token.
// @Description
// @Description  **Akses:** Public (tidak memerlukan autentikasi).
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get data from body", error: "Key: 'Email' Error:..."`
// @Description  - `400` Email tidak terdaftar -> `message: "failed to login user", error: "user not found"`
// @Description  - `400` Password salah -> `message: "failed to login user", error: "crypto/bcrypt: ..."`
// @Description  - `500` Kesalahan internal server -> `message: "failed to login user", error: "Internal Error"`
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      userDto.UserLoginRequest  true  "Payload Login"
// @Success      200  {object}  utils.Response[dto.TokenResponse,any]
// @Failure      400  {object}  swagger.ErrLoginFailed
// @Failure      500  {object}  swagger.ErrLoginInternalServer
// @Router       /api/auth/login [post]
func (c *authController) Login(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var req userDto.UserLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.authService.Login(ctx, req)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(userDto.MESSAGE_FAILED_LOGIN_USER, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(userDto.MESSAGE_FAILED_LOGIN_USER, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(userDto.MESSAGE_SUCCESS_LOGIN_USER, result, path)
	ctx.JSON(http.StatusOK, res)
}

// Logout godoc
// @Summary      User Logout
// @Description  Menghapus session user dan menghapus semua refresh token milik user dari database.
// @Description
// @Description  **Akses:** Authenticated User (memerlukan Bearer Token di header Authorization).
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `400` Logout gagal (non-internal error) -> `message: "failed logout", error: "..."`
// @Description  - `500` Kesalahan internal server -> `message: "failed logout", error: "Internal Error"`
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  utils.Response[any,any]
// @Failure      400  {object}  swagger.ErrLogoutFailed
// @Failure      401  {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      500  {object}  swagger.ErrLogoutFailed
// @Router       /api/auth/logout [post]
func (c *authController) Logout(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)
	path := ctx.Request.URL.Path
	if err := c.authService.Logout(ctx, userId); err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_LOGOUT, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_LOGOUT, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_LOGOUT, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}

// RefreshToken godoc
// @Summary      Refresh Access Token
// @Description  Mendapatkan access token baru menggunakan refresh token yang masih valid.
// @Description  Refresh token lama akan dihapus dan digantikan dengan refresh token baru.
// @Description
// @Description  **Akses:** Public (tidak memerlukan autentikasi, cukup dengan refresh token).
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get data from body", error: "Key: 'RefreshToken' Error:..."`
// @Description  - `400` Refresh token tidak ditemukan di database -> `message: "failed refresh token", error: "refresh token not found"`
// @Description  - `401` Refresh token sudah kedaluwarsa -> `message: "failed refresh token", error: "refresh token expired"`
// @Description  - `500` Kesalahan internal server -> `message: "failed refresh token", error: "Internal Error"`
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RefreshTokenRequest  true  "Payload Refresh Token"
// @Success      200  {object}  utils.Response[dto.TokenResponse,any]
// @Failure      400  {object}  swagger.ErrRefreshTokenNotFound
// @Failure      401  {object}  swagger.ErrRefreshTokenExpired
// @Failure      500  {object}  swagger.ErrRefreshTokenInternalServer
// @Router       /api/auth/refresh-token [post]
func (c *authController) RefreshToken(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.authService.RefreshToken(ctx, req)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REFRESH_TOKEN, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		if errors.Is(err, dto.ErrRefreshTokenExpired) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REFRESH_TOKEN, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REFRESH_TOKEN, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REFRESH_TOKEN, result, path)
	ctx.JSON(http.StatusOK, res)
}

// ResetPassword godoc
// @Summary      Reset Password User
// @Description  Mengubah password user berdasarkan email. Hanya pemilik akun atau Super Admin yang diizinkan.
// @Description
// @Description  **Akses:** Authenticated User — hanya bisa reset password milik sendiri. Super Admin dapat reset password siapa saja.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get data from body", error: "Key: 'Email' Error:..."`
// @Description  - `400` Email tidak terdaftar -> `message: "failed send password reset", error: "user not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `401` Bukan pemilik akun dan bukan Super Admin -> `message: "User unauthorized", error: "You are not authorized to reset this password"`
// @Description  - `500` Kesalahan internal server -> `message: "failed send password reset", error: "Internal Error"`
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      dto.ResetPasswordRequest  true  "Payload Reset Password"
// @Success      200  {object}  utils.Response[any,any]
// @Failure      400  {object}  swagger.ErrResetPasswordFailed
// @Failure      401  {object}  swagger.ErrUnauthorizedResetPassword
// @Failure      500  {object}  swagger.ErrResetPasswordInternalServer
// @Router       /api/auth/reset-password [post]
func (c *authController) ResetPassword(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var req dto.ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userEmail := ctx.MustGet("user_email").(string)
	userRole := ctx.MustGet("role_name").(string)
	if userEmail != req.Email && userRole != constants.ROLE_SUPER_ADMIN {
		res := utils.BuildResponseFailed("User unauthorized", "You are not authorized to reset this password", nil, path)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}

	if err := c.authService.ResetPassword(ctx, req); err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_SEND_PASSWORD_RESET, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_SEND_PASSWORD_RESET, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_RESET_PASSWORD, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}
