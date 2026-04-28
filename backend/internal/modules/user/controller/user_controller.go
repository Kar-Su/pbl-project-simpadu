package controller

import (
	"errors"
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
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	userValidation := validation.NewUserValidation()
	return &userController{
		userService:    userServ,
		roleService:    roleService,
		userValidation: userValidation,
		db:             db,
	}
}

func (c *userController) GetUser(ctx *gin.Context) {
	userId := ctx.Param("id")

	result, err := c.userService.GetUserByID(ctx.Request.Context(), uuid.MustParse(userId))
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
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_USER, result)
	ctx.JSON(http.StatusOK, res)
}

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
			res := utils.BuildResponseFailed(err.Error(), err.Error(), nil)
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
			res := utils.BuildResponseFailed(err.Error(), err.Error(), nil)
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

func (c *userController) GetUserByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	result, err := c.userService.GetUserByEmail(ctx.Request.Context(), email)
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
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_USER, result)
	ctx.JSON(http.StatusOK, res)
}

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
			res := utils.BuildResponseFailed(err.Error(), err.Error(), nil)
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
			res := utils.BuildResponseFailed(err.Error(), err.Error(), nil)
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
			res := utils.BuildResponseFailed(err.Error(), err.Error(), nil)
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

func (c *userController) RegisterAdmin(ctx *gin.Context) {
	var reqBody dto.UserAdminCreateRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	data, err := c.userService.CreateAdmin(ctx.Request.Context(), reqBody)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REGISTER_USER, data)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) RegisterNonAdmin(ctx *gin.Context) {
	var reqBody dto.UserNonAdminCreateRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		res := utils.BuildResponseFailed(constants.MESAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	data, err := c.userService.CreateNonAdmin(ctx.Request.Context(), reqBody)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REGISTER_USER, data)
	ctx.JSON(http.StatusOK, res)
}

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

func (c *userController) DeleteAdmin(ctx *gin.Context) {
	userId := ctx.Param("id")
	if err := c.userService.DeleteAdmin(ctx.Request.Context(), uuid.MustParse(userId)); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_USER, nil)
	ctx.JSON(http.StatusOK, res)
}

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
			res := utils.BuildResponseFailed(err.Error(), err.Error(), nil)
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
