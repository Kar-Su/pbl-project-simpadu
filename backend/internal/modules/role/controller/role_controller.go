package controller

import (
	"errors"
	"log"
	"net/http"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/modules/role/dto"
	"web-hosting/internal/modules/role/service"
	"web-hosting/internal/package/constants"
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
// @Description  Akses: Khusus Super Admin.
// @Tags         role (super)
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      dto.RoleCreateRequest  true  "Payload Create Role"
// @Success      200  {object}  utils.Response{data=entities.Role}
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /api/super/role [post]
func (c *roleController) Create(ctx *gin.Context) {
	var req dto.RoleCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_REQUEST_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	role, err := c.roleService.Create(ctx.Request.Context(), req)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_ROLE, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_ROLE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_ROLE, role)
	ctx.JSON(http.StatusOK, res)
}

// UpdateRole godoc
// @Summary      Update Role
// @Description  Mengubah data role berdasarkan nama role yang ada di URI.
// @Description  Akses: Khusus Super Admin.
// @Tags         role (super)
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        role_name  path      string                 true  "Nama Role Saat Ini"
// @Param        request    body      dto.RoleUpdateRequest  true  "Payload Update Role"
// @Success      200  {object}  utils.Response{data=entities.Role}
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /api/super/role/{role_name} [put]
func (c *roleController) Update(ctx *gin.Context) {
	var RoleNameURI dto.RoleNameURI
	if err := ctx.ShouldBindUri(&RoleNameURI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_VALIDATE_ROLE_URI, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var req dto.RoleUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_REQUEST_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	roleId, err := c.roleService.GetRoleIdByRoleName(ctx.Request.Context(), RoleNameURI.RoleName)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_ROLE, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_ROLE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	updatedRole, err := c.roleService.Update(ctx.Request.Context(), req, roleId)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_ROLE, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_ROLE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_ROLE, updatedRole)
	ctx.JSON(http.StatusOK, res)
}

// DeleteRole godoc
// @Summary      Hapus Role
// @Description  Menghapus role dari sistem berdasarkan nama role.
// @Description  Akses: Khusus Super Admin.
// @Tags         role (super)
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        role_name  path      string  true  "Nama Role yang akan dihapus"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /api/super/role/{role_name} [delete]
func (c *roleController) Delete(ctx *gin.Context) {
	var RoleNameURI dto.RoleNameURI
	log.Printf("roleName: %s\n", RoleNameURI.RoleName)
	if err := ctx.ShouldBindUri(&RoleNameURI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_VALIDATE_ROLE_URI, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	roleName := RoleNameURI.RoleName
	roleId, err := c.roleService.GetRoleIdByRoleName(ctx.Request.Context(), roleName)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_ROLE, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_ROLE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	if err := c.roleService.Delete(ctx.Request.Context(), roleId); err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_ROLE, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_ROLE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_ROLE, nil)
	ctx.JSON(http.StatusOK, res)
}

func (c *roleController) GetAllRole(ctx *gin.Context) {
	roles, err := c.roleService.GetAllRole(ctx.Request.Context())
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ROLE, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ROLE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_ROLE, roles)
	ctx.JSON(http.StatusOK, res)
}
