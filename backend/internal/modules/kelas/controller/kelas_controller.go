package controller

import (
	"errors"
	"net/http"
	"web-hosting/internal/modules/kelas/dto"
	"web-hosting/internal/modules/kelas/service"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type (
	KelasController interface {
		CreateKelas(ctx *gin.Context)
		UpdateKelas(ctx *gin.Context)
		DeleteKelas(ctx *gin.Context)
		GetKelasByID(ctx *gin.Context)
		GetKelasByProdiName(ctx *gin.Context)
	}

	kelasController struct {
		db           *gorm.DB
		kelasService service.KelasService
	}
)

func NewKelasController(injector do.Injector, db *gorm.DB, kelasService service.KelasService) KelasController {
	return &kelasController{
		db:           db,
		kelasService: kelasService,
	}
}

func (c *kelasController) CreateKelas(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var req dto.KelasCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.kelasService.CreateKelas(ctx.Request.Context(), req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_KELAS, any(nil), path)
	ctx.JSON(http.StatusCreated, res)
}

func (c *kelasController) UpdateKelas(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.KelasIdURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var req dto.KelasUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	data, err := c.kelasService.UpdateKelas(ctx.Request.Context(), URI.KelasID, req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_KELAS, data, path)
	ctx.JSON(http.StatusOK, res)
}

func (c *kelasController) DeleteKelas(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.KelasIdURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.kelasService.DeleteKelas(ctx.Request.Context(), URI.KelasID)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_KELAS, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}

func (c *kelasController) GetKelasByID(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.KelasIdURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	data, err := c.kelasService.GetKelasByID(ctx.Request.Context(), URI.KelasID)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_KELAS, data, path)
	ctx.JSON(http.StatusOK, res)
}

func (c *kelasController) GetKelasByProdiName(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.ProdiNameURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	data, err := c.kelasService.GetKelasByProdiName(ctx.Request.Context(), URI.Name)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_KELAS, data, path)
	ctx.JSON(http.StatusOK, res)
}
