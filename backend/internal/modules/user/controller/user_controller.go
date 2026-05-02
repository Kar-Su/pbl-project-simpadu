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
// @Summary      Get Current User
// @Description  Mengambil data profil user yang sedang login berdasarkan token.
// @Tags         user
// @Security     ApiKeyAuth
// @Success      200  {object}  utils.Response{data=dto.UserResponse}
// @Failure      400  {object}  utils.ResponseErr
// @Failure      500  {object}  utils.ResponseErr
// @Router       /api/me [get]
func (c *userController) Me(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)

	result, err := c.userService.GetUserByID(ctx.Request.Context(), uuid.MustParse(userId))
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}

		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_USER, result)
	ctx.JSON(http.StatusOK, res)
}

// GetUser godoc
// @Summary      Get User By ID
// @Description  Mengambil data user berdasarkan UUID.
// @Description  Akses: Super Admin.
// @Tags         user (super)
// @Security     ApiKeyAuth
// @Param        id   path      string  true  "User UUID"
// @Success      200  {object}  utils.Response{data=dto.UserResponse}
// @Failure      400  {object}  utils.ResponseErr
// @Failure      500  {object}  utils.ResponseErr
// @Router       /api/user/{id} [get]
func (c *userController) GetUser(ctx *gin.Context) {
	userId := ctx.Param("id")

	result, err := c.userService.GetUserByID(ctx.Request.Context(), uuid.MustParse(userId))
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}

		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_USER, result)
	ctx.JSON(http.StatusOK, res)
}

// GetUserNonAdmin godoc
// @Summary      Sync/Get User Non-Admin
// @Description  Mendapatkan detail user berdasarkan role dan ID detail.
// @Tags         user
// @Security     ApiKeyAuth
// @Param        role_name  path      string  true  "Role Name"
// @Param        detail_id  path      int     true  "Detail ID"
// @Success      200        {object}  utils.Response{data=dto.UserResponse}
// @Failure      400        {object}  utils.ResponseErr
// @Failure      500        {object}  utils.ResponseErr
// @Router       /api/user/sync/{role_name}/{detail_id} [get]
func (c *userController) GetUserNonAdmin(ctx *gin.Context) {
	var req dto.UserSyncURI
	if err := ctx.ShouldBindUri(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	roleId, err := c.roleService.GetRoleIdByRoleName(ctx.Request.Context(), req.RoleName)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.userService.GetUserByRoleAndDetailID(ctx.Request.Context(), roleId, req.DetailId)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_USER, result)
	ctx.JSON(http.StatusOK, res)
}

// GetUserByEmail godoc
// @Summary      Get User By Email
// @Description  Mencari user spesifik menggunakan email.
// @Tags         user
// @Security     ApiKeyAuth
// @Accept       json
// @Param        request  body      dto.UserEmailRequest  true  "Email Payload"
// @Success      200      {object}  utils.Response{data=dto.UserResponse}
// @Failure      400      {object}  utils.ResponseErr
// @Failure      500      {object}  utils.ResponseErr
// @Router       /api/user/email/ [get]
func (c *userController) GetUserByEmail(ctx *gin.Context) {
	var req dto.UserEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	email := req.Email
	result, err := c.userService.GetUserByEmail(ctx.Request.Context(), email)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}

		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_USER, result)
	ctx.JSON(http.StatusOK, res)
}

// GetUserByRole godoc
// @Summary      Get Users By Role Name
// @Description  Mengambil daftar user berdasarkan nama role (misal: mahasiswa, pegawai).
// @Tags         user
// @Security     ApiKeyAuth
// @Param        role_name  path      string  true  "Role Name"
// @Success      200        {object}  utils.Response{data=[]dto.UserResponse}
// @Failure      400        {object}  utils.ResponseErr
// @Failure      500        {object}  utils.ResponseErr
// @Router       /api/user/role/{role_name} [get]
func (c *userController) GetUserByRole(ctx *gin.Context) {
	var req dto.UserRoleURI
	if err := ctx.ShouldBindUri(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	roleId, err := c.roleService.GetRoleIdByRoleName(ctx.Request.Context(), req.RoleName)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.userService.GetUserByRole(ctx.Request.Context(), roleId)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return

	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_USER, result)
	ctx.JSON(http.StatusOK, res)
}

// UpdateAdmin godoc
// @Summary      Update Admin User
// @Description  Mengupdate data admin berdasarkan ID.
// @Description  Akses: Super Admin.
// @Tags         user (super)
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string true  "User UUID"
// @Param        request  body      dto.UserAdminUpdateRequest  true  "Payload Update"
// @Success      200      {object}  utils.Response{data=dto.UserResponse}
// @Failure      400      {object}  utils.ResponseErr
// @Failure      500      {object}  utils.ResponseErr
// @Router       /api/super/user/{id} [put]
func (c *userController) UpdateAdmin(ctx *gin.Context) {
	var reqBody dto.UserAdminUpdateRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userId := ctx.Param("id")
	data, err := c.userService.UpdateAdmin(ctx.Request.Context(), reqBody, uuid.MustParse(userId))
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_USER, data)
	ctx.JSON(http.StatusOK, res)

}

// RegisterAdmin godoc
// @Summary      Register Admin User
// @Description  Membuat user admin baru (Pegawai/Admin).
// @Description  Akses: Super Admin.
// @Tags         user (super)
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.UserAdminCreateRequest  true  "Payload Admin"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.ResponseErr
// @Failure      500  {object}  utils.ResponseErr
// @Router       /api/super/user [post]
func (c *userController) RegisterAdmin(ctx *gin.Context) {
	var reqBody dto.UserAdminCreateRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	err := c.userService.CreateAdmin(ctx.Request.Context(), reqBody)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REGISTER_USER, nil)
	ctx.JSON(http.StatusOK, res)
}

// RegisterNonAdmin godoc
// @Summary      Register Non-Admin User
// @Description  Membuat user non-admin (Mahasiswa/Pegawai).
// @Description  Akses: Super Admin, Admin Mahasiswa, Admin Pegawai.
// @Tags         user
// @Security     ApiKeyAuth
// @Accept       json
// @Param        request  body      dto.UserNonAdminCreateRequest  true  "Payload Non-Admin"
// @Success      200      {object}  utils.Response
// @Failure      400      {object}  utils.ResponseErr
// @Failure      500      {object}  utils.ResponseErr
//
// @Router       /api/user [post]
func (c *userController) RegisterNonAdmin(ctx *gin.Context) {
	var reqBody dto.UserNonAdminCreateRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	err := c.userService.CreateNonAdmin(ctx.Request.Context(), reqBody)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REGISTER_USER, nil)
	ctx.JSON(http.StatusOK, res)
}

// UpdateNonAdmin godoc
// @Summary      Update User Non-Admin
// @Description  Update profil user. Hanya bisa dilakukan oleh pemilik akun atau Super Admin.
// @Tags         user
// @Security     ApiKeyAuth
// @Accept       json
// @Param        role_name  path      string                         true  "Role Name"
// @Param        detail_id  path      int                            true  "Detail ID"
// @Param        request    body      dto.UserNonAdminUpdateRequest  true  "Update Payload"
// @Success      200        {object}  utils.Response{data=dto.UserResponse}
// @Failure      400        {object}  utils.ResponseErr
// @Failure      500        {object}  utils.ResponseErr
// @Router       /api/user/sync/{role_name}/{detail_id} [put]
func (c *userController) UpdateNonAdmin(ctx *gin.Context) {
	var reqUri dto.UserSyncURI
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var reqBody dto.UserNonAdminUpdateRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userRoleName := ctx.MustGet("role_name").(string)
	userDetailId := ctx.MustGet("detail_id").(uint)

	if userDetailId != reqUri.DetailId && userRoleName != constants.ROLE_SUPER_ADMIN {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, "Unauthorized", nil)
		log.Println(userRoleName)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}

	roleId, err := c.roleService.GetRoleIdByRoleName(ctx.Request.Context(), reqUri.RoleName)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(err.Error(), err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	data, err := c.userService.UpdateNonAdmin(ctx.Request.Context(), reqBody, roleId, reqUri.DetailId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_USER, data)
	ctx.JSON(http.StatusOK, res)
}

// DeleteAdmin godoc
// @Summary      Delete Admin User
// @Description  Menghapus admin secara permanen.
// @Description  Akses: Super Admin.
// @Tags         user (super)
// @Security     ApiKeyAuth
// @Param        id   path      string  true  "User UUID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.ResponseErr
// @Failure      500  {object}  utils.ResponseErr
// @Router       /api/super/user/{id} [delete]
func (c *userController) DeleteAdmin(ctx *gin.Context) {
	userId := ctx.Param("id")
	if err := c.userService.DeleteAdmin(ctx.Request.Context(), uuid.MustParse(userId)); err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_USER, nil)
	ctx.JSON(http.StatusOK, res)
}

// DeleteNonAdmin godoc
// @Summary      Delete User Non-Admin
// @Description  Menghapus user non-admin.
// @Description  Akses: Super Admin, Admin Pegawai, Admin Mahasiswa.
// @Tags         user
// @Security     ApiKeyAuth
// @Param        role_name  path      string  true  "Role Name"
// @Param        detail_id  path      int     true  "Detail ID"
// @Success      200        {object}  utils.Response
// @Failure      400        {object}  utils.ResponseErr
// @Failure      500        {object}  utils.ResponseErr
// @Router       /api/user/sync/{role_name}/{detail_id} [delete]
func (c *userController) DeleteNonAdmin(ctx *gin.Context) {
	var reqUri dto.UserSyncURI
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	roleId, err := c.roleService.GetRoleIdByRoleName(ctx, reqUri.RoleName)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := c.userService.DeleteNonAdmin(ctx.Request.Context(), roleId, reqUri.DetailId); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_USER, nil)
	ctx.JSON(http.StatusOK, res)
}
