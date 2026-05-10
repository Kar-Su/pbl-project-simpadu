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
	KelasMahasiswaController interface {
		AssignMahasiswaToKelas(ctx *gin.Context)
		RemoveMahasiswaFromKelas(ctx *gin.Context)
		GetAllKelasMahasiswa(ctx *gin.Context)
		GetMahasiswaByKelas(ctx *gin.Context)
	}
	kelasMahasiswaController struct {
		db           *gorm.DB
		pivotService service.KelasMahasiswaService
	}
)

func NewKelasMahasiswaController(injector do.Injector, db *gorm.DB, pivotService service.KelasMahasiswaService) KelasMahasiswaController {
	return &kelasMahasiswaController{
		db:           db,
		pivotService: pivotService,
	}
}

func (c *kelasMahasiswaController) AssignMahasiswaToKelas(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var req dto.KelasMahasiswaCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_ASSIGN_MAHASISWA_TO_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	if err := c.pivotService.Create(ctx.Request.Context(), req); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_ASSIGN_MAHASISWA_TO_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_ASSIGN_MAHASISWA_TO_KELAS, any(nil), path)
	ctx.JSON(http.StatusCreated, res)
}

func (c *kelasMahasiswaController) RemoveMahasiswaFromKelas(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.KelasMahasiswaURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REMOVE_MAHASISWA_FROM_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	if err := c.pivotService.Delete(ctx.Request.Context(), URI.MahasiswaID, URI.KelasID); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REMOVE_MAHASISWA_FROM_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REMOVE_MAHASISWA_FROM_KELAS, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}

func (c *kelasMahasiswaController) GetAllKelasMahasiswa(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.MahasiswaIdURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_PIVOT, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	data, err := c.pivotService.GetAllKelasMahasiswa(ctx.Request.Context(), URI.MahasiswaID)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_PIVOT, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DATA_PIVOT, data, path)
	ctx.JSON(http.StatusOK, res)
}

func (c *kelasMahasiswaController) GetMahasiswaByKelas(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.KelasIdURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_PIVOT, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	data, err := c.pivotService.GetMahasiswaByKelasId(ctx.Request.Context(), URI.KelasID)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_PIVOT, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DATA_PIVOT, data, path)
	ctx.JSON(http.StatusOK, res)
}
