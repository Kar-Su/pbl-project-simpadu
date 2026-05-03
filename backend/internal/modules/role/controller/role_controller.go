package controller

import (
	"errors"
	"log"
	"net/http"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/modules/role/dto"
	"web-hosting/internal/modules/role/service"
	"web-hosting/internal/package/constants"
	_ "web-hosting/internal/package/swagger"
	"web-hosting/internal/package/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

var _ = entities.Role{}

type RoleController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetAllRole(ctx *gin.Context)
}

type roleController struct {
	roleService service.RoleService
	db          *gorm.DB
}

func NewRoleController(injector do.Injector, roleService service.RoleService, db *gorm.DB) RoleController {
	return &roleController{
		roleService: roleService,
		db:          db,
	}
}

// CreateRole godoc
// @Summary      Buat Role Baru
// @Description  Menambahkan role baru ke dalam sistem.
// @Description
// @Description  **Akses:** Super Admin.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'RoleName' Error:..."`
// @Description  - `400` Role dengan nama tersebut sudah ada -> `message: "failed to create role", error: "role already exists"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` Role user tidak memiliki akses -> `message: "Role anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to create role", error: "Internal Error"`
// @Tags         role (super)
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      dto.RoleCreateRequest  true  "Payload Create Role"
// @Success      200      {object}  utils.Response[entities.Role,any]
// @Failure      400      {object}  swagger.ErrCreateRoleFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrCreateRoleInternalServer
// @Router       /api/super/role [post]
func (c *roleController) Create(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var req dto.RoleCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_REQUEST_BODY, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	role, err := c.roleService.Create(ctx.Request.Context(), req)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_ROLE, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_ROLE, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_ROLE, role, path)
	ctx.JSON(http.StatusOK, res)
}

// UpdateRole godoc
// @Summary      Update Role
// @Description  Mengubah nama role berdasarkan nama role yang sudah ada.
// @Description
// @Description  **Akses:** Super Admin.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter URI tidak valid -> `message: "failed to validate role uri", error: "Key: 'RoleName' Error:..."`
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'RoleName' Error:..."`
// @Description  - `400` Role dengan nama tersebut tidak ditemukan -> `message: "failed to update role", error: "role not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` Role user tidak memiliki akses -> `message: "Role anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to update role", error: "Internal Error"`
// @Tags         role (super)
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        role_name  path      string                 true  "Nama Role yang Akan Diubah"  example(mahasiswa)
// @Param        request    body      dto.RoleUpdateRequest  true  "Payload Update Role"
// @Success      200        {object}  utils.Response[entities.Role,any]
// @Failure      400        {object}  swagger.ErrUpdateRoleFailed
// @Failure      401        {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403        {object}  swagger.ErrForbiddenAccess
// @Failure      500        {object}  swagger.ErrUpdateRoleInternalServer
// @Router       /api/super/role/{role_name} [put]
func (c *roleController) Update(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var RoleNameURI dto.RoleNameURI
	if err := ctx.ShouldBindUri(&RoleNameURI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_VALIDATE_ROLE_URI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var req dto.RoleUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_REQUEST_BODY, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	roleId, err := c.roleService.GetRoleIdByRoleName(ctx.Request.Context(), RoleNameURI.RoleName)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_ROLE, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_ROLE, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	updatedRole, err := c.roleService.Update(ctx.Request.Context(), req, roleId)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_ROLE, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_ROLE, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_ROLE, updatedRole, path)
	ctx.JSON(http.StatusOK, res)
}

// DeleteRole godoc
// @Summary      Hapus Role
// @Description  Menghapus role dari sistem secara permanen berdasarkan nama role.
// @Description
// @Description  **Akses:** Khusus Super Admin.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter URI tidak valid -> `message: "failed to validate role uri", error: "Key: 'RoleName' Error:..."`
// @Description  - `400` Role dengan nama tersebut tidak ditemukan -> `message: "failed to delete role", error: "role not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` Role user tidak memiliki akses -> `message: "Role anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to delete role", error: "Internal Error"`
// @Tags         role (super)
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        role_name  path      string  true  "Nama Role yang Akan Dihapus"  example(mahasiswa)
// @Success      200        {object}  utils.Response[any,any]
// @Failure      400        {object}  swagger.ErrDeleteRoleFailed
// @Failure      401        {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403        {object}  swagger.ErrForbiddenAccess
// @Failure      500        {object}  swagger.ErrDeleteRoleInternalServer
// @Router       /api/super/role/{role_name} [delete]
func (c *roleController) Delete(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var RoleNameURI dto.RoleNameURI
	log.Printf("roleName: %s\n", RoleNameURI.RoleName)
	if err := ctx.ShouldBindUri(&RoleNameURI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_VALIDATE_ROLE_URI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	roleName := RoleNameURI.RoleName
	roleId, err := c.roleService.GetRoleIdByRoleName(ctx.Request.Context(), roleName)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_ROLE, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_ROLE, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	if err := c.roleService.Delete(ctx.Request.Context(), roleId); err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_ROLE, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_ROLE, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_ROLE, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}

// GetAllRole godoc
// @Summary      Ambil Semua Role
// @Description  Mengambil seluruh daftar role yang tersedia di sistem.
// @Description
// @Description  **Akses:** Super Admin, Admin Akademik, Admin Keuangan, Admin Mahasiswa.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` Role user tidak memiliki akses -> `message: "Role anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to get role", error: "Internal Error"`
// @Tags         role
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  utils.Response[[]entities.Role,any]
// @Failure      401  {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403  {object}  swagger.ErrForbiddenAccess
// @Failure      500  {object}  swagger.ErrGetRoleFailed
// @Router       /api/role [get]
func (c *roleController) GetAllRole(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	roles, err := c.roleService.GetAllRole(ctx.Request.Context())
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ROLE, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ROLE, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_ROLE, roles, path)
	ctx.JSON(http.StatusOK, res)
}
