package controller

import (
	"errors"
	"log"
	"net/http"
	roleService "web-hosting/internal/modules/role/service"
	"web-hosting/internal/modules/user/dto"
	"web-hosting/internal/modules/user/service"
	"web-hosting/internal/modules/user/validation"
	"web-hosting/internal/package/constants"
	_ "web-hosting/internal/package/swagger"
	"web-hosting/internal/package/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type UserController interface {
	Me(ctx *gin.Context)
	GetUser(ctx *gin.Context)
	GetUserNonAdmin(ctx *gin.Context)
	GetUserByEmail(ctx *gin.Context)
	GetUserByRole(ctx *gin.Context)
	RegisterAdmin(ctx *gin.Context)
	RegisterNonAdmin(ctx *gin.Context)
	UpdateAdmin(ctx *gin.Context)
	UpdateNonAdmin(ctx *gin.Context)
	DeleteAdmin(ctx *gin.Context)
	DeleteNonAdmin(ctx *gin.Context)
}

type userController struct {
	userService    service.UserService
	roleService    roleService.RoleService
	userValidation *validation.UserValidation
	db             *gorm.DB
}

func NewUserController(injector do.Injector, userServ service.UserService, roleService roleService.RoleService) UserController {
	// db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB_TEST)
	userValidation := validation.NewUserValidation()
	return &userController{
		userService:    userServ,
		roleService:    roleService,
		userValidation: userValidation,
		db:             db,
	}
}

// Me godoc
// @Summary      Get Current User (Profil Saya)
// @Description  Mengambil data profil user yang sedang login berdasarkan JWT token.
// @Description
// @Description  **Akses:** Semua user yang sudah login (Authenticated User).
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `400` User tidak ditemukan -> `message: "failed to get user", error: "user not found"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to get user", error: "Internal Error"`
// @Tags         user
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  utils.Response[dto.UserResponse,any]
// @Failure      400  {object}  swagger.ErrGetUserFailed
// @Failure      401  {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      500  {object}  swagger.ErrGetUserInternalServer
// @Router       /api/me [get]
func (c *userController) Me(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)
	path := ctx.Request.URL.Path

	result, err := c.userService.GetUserByID(ctx.Request.Context(), uuid.MustParse(userId))
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}

		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_USER, result, path)
	ctx.JSON(http.StatusOK, res)
}

// GetUser godoc
// @Summary      Get User By ID
// @Description  Mengambil data lengkap seorang user berdasarkan UUID-nya.
// @Description
// @Description  **Akses:** Khusus Super Admin.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` User dengan ID tersebut tidak ditemukan -> `message: "failed to get user", error: "user not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` Role user tidak memiliki akses -> `message: "Role anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to get user", error: "Internal Error"`
// @Tags         user (super)
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      string  true  "UUID User"  example(019748ae-beef-7abc-b123-abcdef012345)
// @Success      200  {object}  utils.Response[dto.UserResponse,any]
// @Failure      400  {object}  swagger.ErrGetUserFailed
// @Failure      401  {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403  {object}  swagger.ErrForbiddenAccess
// @Failure      500  {object}  swagger.ErrGetUserInternalServer
// @Router       /api/user/{id} [get]
func (c *userController) GetUser(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	userId := ctx.Param("id")

	result, err := c.userService.GetUserByID(ctx.Request.Context(), uuid.MustParse(userId))
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}

		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_USER, result, path)
	ctx.JSON(http.StatusOK, res)
}

// GetUserNonAdmin godoc
// @Summary      Get / Sync User Non-Admin
// @Description  Mendapatkan data user non-admin berdasarkan role_name dan detail_id (misalnya NIM mahasiswa atau NIP pegawai).
// @Description  Endpoint ini digunakan untuk sinkronisasi data antara sistem eksternal dengan sistem ini.
// @Description
// @Description  **Akses:** Semua user yang sudah login (Authenticated User).
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter URI tidak valid (role_name atau detail_id salah format) -> `message: "bad request", error: "Key: 'RoleName' Error:..."`
// @Description  - `400` Role tidak ditemukan -> `message: "failed to get user", error: "role not found"`
// @Description  - `400` User tidak ditemukan -> `message: "failed to get user", error: "user not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to get user", error: "Internal Error"`
// @Tags         user
// @Produce      json
// @Security     ApiKeyAuth
// @Param        role_name  path      string  true  "Nama Role Non-Admin"  example(mahasiswa)
// @Param        detail_id  path      int     true  "ID Detail (misal NIM/NIP)"  example(10)
// @Success      200  {object}  utils.Response[dto.UserResponse,any]
// @Failure      400  {object}  swagger.ErrGetUserFailed
// @Failure      401  {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      500  {object}  swagger.ErrGetUserInternalServer
// @Router       /api/user/sync/{role_name}/{detail_id} [get]
func (c *userController) GetUserNonAdmin(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var req dto.UserSyncURI
	if err := ctx.ShouldBindUri(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, err.Error(), nil, path)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	roleId, err := c.roleService.GetRoleIdByRoleName(ctx.Request.Context(), req.RoleName)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.userService.GetUserByRoleAndDetailID(ctx.Request.Context(), roleId, req.DetailId)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_USER, result)
	ctx.JSON(http.StatusOK, res)
}

// GetUserByEmail godoc
// @Summary      Get User By Email
// @Description  Mencari dan mengambil data user spesifik menggunakan alamat email.
// @Description
// @Description  **Akses:** Semua user yang sudah login (Authenticated User).
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong atau email tidak valid -> `message: "failed to get data from body", error: "Key: 'Email' Error:..."`
// @Description  - `400` User dengan email tersebut tidak ditemukan -> `message: "failed to get user", error: "user not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to get user", error: "Internal Error"`
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        email  query      dto.UserEmailRequest  true  "Payload Email User" example(rezi@example.com)
// @Success      200  {object}  utils.Response[dto.UserResponse,any]
// @Failure      400  {object}  swagger.ErrGetUserFailed
// @Failure      401  {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      500  {object}  swagger.ErrGetUserInternalServer
// @Router       /api/user/email/ [get]
func (c *userController) GetUserByEmail(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var req dto.UserEmailRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	email := req.Email
	result, err := c.userService.GetUserByEmail(ctx.Request.Context(), email)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}

		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_USER, result, path)
	ctx.JSON(http.StatusOK, res)
}

// GetUserByRole godoc
// @Summary      Get Users By Role Name
// @Description  Mengambil daftar semua user yang memiliki role tertentu.
// @Description
// @Description  **Akses:** Semua user yang sudah login (Authenticated User).
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter URI tidak valid -> `message: "bad request", error: "Key: 'RoleName' Error:..."`
// @Description  - `400` Role tidak ditemukan -> `message: "failed to get user", error: "role not found"`
// @Description  - `400` Gagal mengambil daftar user -> `message: "failed to get user", error: "..."`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to get user", error: "Internal Error"`
// @Tags         user
// @Produce      json
// @Security     ApiKeyAuth
// @Param        role_name  path      string  true  "Nama Role"  example(mahasiswa)
// @Success      200        {object}  utils.Response[[]dto.UserResponse,any]
// @Failure      400        {object}  swagger.ErrGetListUserFailed
// @Failure      401        {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      500        {object}  swagger.ErrGetUserInternalServer
// @Router       /api/user/role/{role_name} [get]
func (c *userController) GetUserByRole(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var req dto.UserRoleURI
	if err := ctx.ShouldBindUri(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	roleId, err := c.roleService.GetRoleIdByRoleName(ctx.Request.Context(), req.RoleName)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.userService.GetUserByRole(ctx.Request.Context(), roleId)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return

	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_USER, result, path)
	ctx.JSON(http.StatusOK, res)
}

// UpdateAdmin godoc
// @Summary      Update Admin User
// @Description  Mengupdate data user admin (nama, email, password, role, image) berdasarkan UUID.
// @Description  Semua field bersifat opsional — hanya field yang diisi yang akan diupdate.
// @Description
// @Description  **Akses:** Khusus Super Admin.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid -> `message: "failed to get data from body", error: "Key: 'Name' Error:..."`
// @Description  - `400` User tidak ditemukan -> `message: "failed to update user", error: "user not found"`
// @Description  - `400` Email sudah digunakan user lain -> `message: "failed to update user", error: "email already exists"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` Role user tidak memiliki akses -> `message: "Role anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to update user", error: "Internal Error"`
// @Tags         user (super)
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id       path      string                      true  "UUID User Admin"  example(019748ae-beef-7abc-b123-abcdef012345)
// @Param        request  body      swagger.UserAdminUpdateRequest  true  "Payload Update Admin"
// @Success      200  {object}  utils.Response[dto.UserResponse,any]
// @Failure      400  {object}  swagger.ErrUpdateUserFailed
// @Failure      401  {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403  {object}  swagger.ErrForbiddenAccess
// @Failure      500  {object}  swagger.ErrUpdateUserInternalServer
// @Router       /api/super/user/{id} [put]
func (c *userController) UpdateAdmin(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var reqBody dto.UserAdminUpdateRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userId := ctx.Param("id")
	data, err := c.userService.UpdateAdmin(ctx.Request.Context(), reqBody, uuid.MustParse(userId))
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_USER, data, path)
	ctx.JSON(http.StatusOK, res)

}

// RegisterAdmin godoc
// @Summary      Register Admin User
// @Description  Membuat akun user admin baru (Admin Akademik, Admin Keuangan, Admin Mahasiswa, Admin Pegawai, dsb).
// @Description
// @Description  **Akses:** Khusus Super Admin.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get data from body", error: "Key: 'Email' Error:..."`
// @Description  - `400` Email sudah terdaftar -> `message: "failed to register user", error: "email already exists"`
// @Description  - `400` Role tidak valid -> `message: "failed to register user", error: "role not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` Role user tidak memiliki akses -> `message: "Role anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to register user", error: "Internal Error"`
// @Tags         user (super)
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      swagger.UserAdminCreateRequest  true  "Payload Registrasi Admin"
// @Success      200  {object}  utils.Response[any,any]
// @Failure      400  {object}  swagger.ErrRegisterUserFailed
// @Failure      401  {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403  {object}  swagger.ErrForbiddenAccess
// @Failure      500  {object}  swagger.ErrRegisterUserInternalServer
// @Router       /api/super/user [post]
func (c *userController) RegisterAdmin(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var reqBody dto.UserAdminCreateRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	err := c.userService.CreateAdmin(ctx.Request.Context(), reqBody)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REGISTER_USER, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}

// RegisterNonAdmin godoc
// @Summary      Register Non-Admin User
// @Description  Membuat akun user non-admin baru (Mahasiswa, Dosen, Pegawai, dsb).
// @Description  Field `role_name` harus merupakan role non-admin (bukan super-admin, admin-*).
// @Description
// @Description  **Akses:** Super Admin, Admin Mahasiswa, Admin Pegawai.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get data from body", error: "Key: 'Email' Error:..."`
// @Description  - `400` Email sudah terdaftar -> `message: "failed to register user", error: "email already exists"`
// @Description  - `400` Role tidak valid atau role adalah role admin -> `message: "failed to register user", error: "invalid not admin role"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` Role user tidak memiliki akses -> `message: "Role anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to register user", error: "Internal Error"`
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      swagger.UserNonAdminCreateRequest  true  "Payload Registrasi Non-Admin"
// @Success      200  {object}  utils.Response[any,any]
// @Failure      400  {object}  swagger.ErrRegisterUserFailed
// @Failure      401  {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403  {object}  swagger.ErrForbiddenAccess
// @Failure      500  {object}  swagger.ErrRegisterUserInternalServer
// @Router       /api/user [post]
func (c *userController) RegisterNonAdmin(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var reqBody dto.UserNonAdminCreateRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	err := c.userService.CreateNonAdmin(ctx.Request.Context(), reqBody)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REGISTER_USER, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}

// UpdateNonAdmin godoc
// @Summary      Update User Non-Admin
// @Description  Mengupdate data profil user non-admin berdasarkan role dan detail_id.
// @Description  Semua field body bersifat opsional — hanya field yang diisi yang akan diupdate.
// @Description
// @Description  **Akses:** Pemilik akun sendiri (detail_id harus cocok) atau Super Admin.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter URI tidak valid -> `message: "bad request", error: "Key: 'RoleName' Error:..."`
// @Description  - `400` Body tidak valid -> `message: "failed to get data from body", error: "Key: 'Email' Error:..."`
// @Description  - `400` Role tidak ditemukan -> `message: "failed to get user", error: "role not found"`
// @Description  - `400` Gagal update user -> `message: "failed to update user", error: "user not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `401` Bukan pemilik akun dan bukan Super Admin -> `message: "failed to update user", error: "Unauthorized"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to update user", error: "Internal Error"`
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        role_name  path      string                         true  "Nama Role Non-Admin"  example(mahasiswa)
// @Param        detail_id  path      int                            true  "ID Detail (misal NIM/NIP)"  example(10)
// @Param        request    body      swagger.UserNonAdminUpdateRequest  true  "Payload Update Non-Admin"
// @Success      200  {object}  utils.Response[dto.UserResponse,any]
// @Failure      400  {object}  swagger.ErrUpdateUserFailed
// @Failure      401  {object}  swagger.ErrUnauthorizedUpdateNonAdmin
// @Failure      500  {object}  swagger.ErrUpdateUserInternalServer
// @Router       /api/user/sync/{role_name}/{detail_id} [put]
func (c *userController) UpdateNonAdmin(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var reqUri dto.UserSyncURI
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var reqBody dto.UserNonAdminUpdateRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userRoleName := ctx.MustGet("role_name").(string)
	userDetailId := ctx.MustGet("detail_id").(uint)

	if userDetailId != reqUri.DetailId && userRoleName != constants.ROLE_SUPER_ADMIN {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, "Unauthorized", nil, path)
		log.Println(userRoleName)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}

	roleId, err := c.roleService.GetRoleIdByRoleName(ctx.Request.Context(), reqUri.RoleName)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(err.Error(), err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	data, err := c.userService.UpdateNonAdmin(ctx.Request.Context(), reqBody, roleId, reqUri.DetailId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_USER, data, path)
	ctx.JSON(http.StatusOK, res)
}

// DeleteAdmin godoc
// @Summary      Delete Admin User
// @Description  Menghapus akun admin secara permanen dari sistem berdasarkan UUID.
// @Description
// @Description  **Akses:** Khusus Super Admin.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` User tidak ditemukan -> `message: "failed to delete user", error: "user not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` Role user tidak memiliki akses -> `message: "Role anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to delete user", error: "Internal Error"`
// @Tags         user (super)
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      string  true  "UUID User Admin"  example(019748ae-beef-7abc-b123-abcdef012345)
// @Success      200  {object}  utils.Response[any,any]
// @Failure      400  {object}  swagger.ErrDeleteUserFailed
// @Failure      401  {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403  {object}  swagger.ErrForbiddenAccess
// @Failure      500  {object}  swagger.ErrDeleteUserInternalServer
// @Router       /api/super/user/{id} [delete]
func (c *userController) DeleteAdmin(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	userId := ctx.Param("id")
	if err := c.userService.DeleteAdmin(ctx.Request.Context(), uuid.MustParse(userId)); err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_USER, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}

// DeleteNonAdmin godoc
// @Summary      Delete User Non-Admin
// @Description  Menghapus akun user non-admin secara permanen berdasarkan role_name dan detail_id.
// @Description
// @Description  **Akses:** Super Admin, Admin Pegawai, Admin Mahasiswa.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter URI tidak valid -> `message: "bad request", error: "Key: 'RoleName' Error:..."`
// @Description  - `400` Role tidak ditemukan -> `message: "failed to get user", error: "role not found"`
// @Description  - `400` User tidak ditemukan -> `message: "failed to delete user", error: "user not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` Role user tidak memiliki akses -> `message: "Role anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to delete user", error: "Internal Error"`
// @Tags         user
// @Produce      json
// @Security     ApiKeyAuth
// @Param        role_name  path      string  true  "Nama Role Non-Admin"  example(mahasiswa)
// @Param        detail_id  path      int     true  "ID Detail (misal NIM/NIP)"  example(10)
// @Success      200  {object}  utils.Response[any,any]
// @Failure      400  {object}  swagger.ErrDeleteUserFailed
// @Failure      401  {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403  {object}  swagger.ErrForbiddenAccess
// @Failure      500  {object}  swagger.ErrDeleteUserInternalServer
// @Router       /api/user/sync/{role_name}/{detail_id} [delete]
func (c *userController) DeleteNonAdmin(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var reqUri dto.UserSyncURI
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	roleId, err := c.roleService.GetRoleIdByRoleName(ctx, reqUri.RoleName)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := c.userService.DeleteNonAdmin(ctx.Request.Context(), roleId, reqUri.DetailId); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_USER, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}
