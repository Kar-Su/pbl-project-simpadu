package controller

import (
	"errors"
	"net/http"
	"web-hosting/internal/modules/mk/dto"
	"web-hosting/internal/modules/mk/service"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/helpers"
	_ "web-hosting/internal/package/swagger"
	"web-hosting/internal/package/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type (
	MkController interface {
		CreateMk(ctx *gin.Context)
		UpdateMk(ctx *gin.Context)
		DeleteMk(ctx *gin.Context)
		GetMk(ctx *gin.Context)
	}

	mkController struct {
		mkService service.MkService
		db        *gorm.DB
	}
)

func NewMkController(injector do.Injector, mkService service.MkService, db *gorm.DB) MkController {
	return &mkController{
		mkService: mkService,
		db:        db,
	}
}

// CreateMataKuliah godoc
// @Summary Create Mata Kuliah Baru
// @Description Menambahkan mata kuliah baru ke sistem
// @Description
// @Description  **Akses:** Admin Akademik
// @Tags mata-kuliah
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.MkCreateRequest true "Mata Kuliah Request"
// @Success      201      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrCreateMkFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrCreateMkInternalServer
// @Router /api/mata-kuliah [post]
func (c *mkController) CreateMk(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var req dto.MkCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_MK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := c.mkService.CreateMk(ctx, req); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_MK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_MK, any(nil), path)
	ctx.JSON(http.StatusCreated, res)
}

// UpdateMataKuliah godoc
// @Summary Update Mata Kuliah berdasarkan kode
// @Description Mengupdate mata kuliah berdasarkan kode di path parameter
// @Description
// @Description  **Akses:** Admin Akademik
// @Tags mata-kuliah
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param kode path string true "Kode Mata Kuliah" example(MK001)
// @Param request body dto.MkUpdateRequest true "Mata Kuliah Request"
// @Success      200      {object}  utils.Response[dto.MkResponse,any]
// @Failure      400      {object}  swagger.ErrUpdateMkFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrUpdateMkInternalServer
// @Router /api/mata-kuliah/{kode} [put]
func (c *mkController) UpdateMk(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var uriParam dto.MkKodeURI
	if err := ctx.ShouldBindUri(&uriParam); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_MK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	kode := helpers.NormalizeString(uriParam.Kode)

	var req dto.MkUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_MK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	mkResponse, err := c.mkService.UpdateMkByKode(ctx, kode, req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_MK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_MK, mkResponse, path)
	ctx.JSON(http.StatusOK, res)
}

// DeleteMataKuliah godoc
// @Summary Delete Mata Kuliah berdasarkan kode
// @Description Menghapus mata kuliah berdasarkan kode di path parameter
// @Description
// @Description  **Akses:** Admin Akademik
// @Tags mata-kuliah
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param kode path string true "Kode Mata Kuliah" example(MK001)
// @Success      200      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrDeleteMkFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrDeleteMkInternalServer
// @Router /api/mata-kuliah/{kode} [delete]
func (c *mkController) DeleteMk(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var uriParam dto.MkKodeURI
	if err := ctx.ShouldBindUri(&uriParam); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_MK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	kode := helpers.NormalizeString(uriParam.Kode)

	if err := c.mkService.DeleteMkByKode(ctx, kode); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_MK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_MK, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}

// GetMataKuliah godoc
// @Summary Get Mata Kuliah
// @Description Mendapatkan data mata kuliah. Jika query id/kode diberikan, return 1 data; tanpa query, return semua dengan pagination
// @Description
// @Description  **Akses:** Logged User
// @Tags mata-kuliah
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query string false "ID Mata Kuliah (UUID)" example(12345678-1234-1234-1234-123456789012)
// @Param kode query string false "Kode Mata Kuliah" example(MK001)
// @Param page query int false "Halaman (default: 1, per halaman: 10)" example(1)
// @Success      200      {object}  utils.Response[dto.MkResponse,any]
// @Success      200      {object}  utils.Response[utils.PaginatedData[[]dto.MkResponse],any]
// @Failure      400      {object}  swagger.ErrUpdateMkFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      500      {object}  swagger.ErrUpdateMkInternalServer
// @Router /api/mata-kuliah [get]
func (c *mkController) GetMk(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var queryParams dto.MkQuery
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_MK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if queryParams.ID != "" && queryParams.Kode != "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_MK, dto.ErrQueryParams.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if queryParams.ID != "" {
		mkId, err := uuid.Parse(queryParams.ID)
		if err != nil {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_MK, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}
		result, err := c.mkService.GetMkById(ctx, mkId)
		if err != nil {
			status := http.StatusBadRequest
			if errors.Is(err, constants.ErrInternalErr) {
				status = http.StatusInternalServerError
			}
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_MK, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(status, res)
			return
		}
		res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_MK, result, path)
		ctx.JSON(http.StatusOK, res)
		return
	}

	if queryParams.Kode != "" {
		queryParams.Kode = helpers.NormalizeString(queryParams.Kode)
		result, err := c.mkService.GetMkByKode(ctx, queryParams.Kode)
		if err != nil {
			status := http.StatusBadRequest
			if errors.Is(err, constants.ErrInternalErr) {
				status = http.StatusInternalServerError
			}
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_MK, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(status, res)
			return
		}
		res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_MK, result, path)
		ctx.JSON(http.StatusOK, res)
		return
	}

	// Pagination untuk get all
	var pageQuery utils.PaginationQuery
	if err := ctx.ShouldBindQuery(&pageQuery); err != nil || pageQuery.Page <= 0 {
		pageQuery.Page = 1
	}

	mks, total, err := c.mkService.GetAllMkPaginated(ctx, pageQuery.Page)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_MK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	paginated := utils.BuildPaginatedResponse(mks, pageQuery.Page, total)
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_MK, paginated, path)
	ctx.JSON(http.StatusOK, res)
}
