package controller

import (
	"errors"
	"net/http"
	"web-hosting/internal/modules/prodi/dto"
	"web-hosting/internal/modules/prodi/service"
	"web-hosting/internal/modules/prodi/validation"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/helpers"
	_ "web-hosting/internal/package/swagger"
	"web-hosting/internal/package/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type (
	ProdiController interface {
		CreateProdi(ctx *gin.Context)
		UpdateProdi(ctx *gin.Context)
		DeleteProdi(ctx *gin.Context)
		GetProdi(ctx *gin.Context)
		GetProdiByJurusan(ctx *gin.Context)
	}

	prodiController struct {
		prodiService   service.ProdiService
		prodiValidator *validation.ProdiValidation
		db             *gorm.DB
	}
)

func NewProdiController(injector do.Injector, prodiService service.ProdiService, db *gorm.DB) ProdiController {
	prodiValidation := validation.NewProdiValidation()
	return &prodiController{
		prodiService:   prodiService,
		prodiValidator: prodiValidation,
		db:             db,
	}
}

// CreateProdi godoc
// @Summary Create Prodi Baru
// @Description Menambahkan prodi baru ke sistem
// @Description
// @Description  **Akses:** Admin Akademik
// @Tags prodi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.ProdiCreateRequest true "Prodi Request"
// @Success      201      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrCreateProdiFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrCreateProdiInternalServer
// @Router /api/prodi [post]
func (c *prodiController) CreateProdi(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var req dto.ProdiCreateRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_PRODI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := c.prodiService.CreateProdi(ctx.Request.Context(), req); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_PRODI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_PRODI, any(nil), path)
	ctx.JSON(http.StatusCreated, res)
}

// UpdateProdi godoc
// @Summary Update Prodi
// @Description Mengupdate prodi berdasarkan nama di path parameter
// @Description
// @Description  **Akses:** Admin Akademik
// @Tags prodi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name path string true "Nama Prodi yang akan diupdate" example(teknik-elektro)
// @Param request body dto.ProdiUpdateRequest true "Prodi Request"
// @Success      200      {object}  utils.Response[dto.ProdiResponse,any]
// @Failure      400      {object}  swagger.ErrUpdateProdiFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrUpdateProdiInternalServer
// @Router /api/prodi/{name} [put]
func (c *prodiController) UpdateProdi(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var uriParam dto.ProdiNameURI
	if err := ctx.ShouldBindUri(&uriParam); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_PRODI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	prodiName := helpers.NormalizeString(uriParam.Name)

	var req dto.ProdiUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_PRODI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	updatedResponse, err := c.prodiService.UpdateProdi(ctx.Request.Context(), prodiName, req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_PRODI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_PRODI, updatedResponse, path)
	ctx.JSON(http.StatusOK, res)
}

// DeleteProdi godoc
// @Summary Delete Prodi
// @Description Menghapus prodi berdasarkan nama di path parameter
// @Description
// @Description  **Akses:** Admin Akademik
// @Tags prodi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name path string true "Nama Prodi yang akan dihapus" example(teknik-elektro)
// @Success      200      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrDeleteProdiFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrDeleteProdiInternalServer
// @Router /api/prodi/{name} [delete]
func (c *prodiController) DeleteProdi(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var uriParam dto.ProdiNameURI
	if err := ctx.ShouldBindUri(&uriParam); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_PRODI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	prodiName := helpers.NormalizeString(uriParam.Name)

	if err := c.prodiService.DeleteProdi(ctx.Request.Context(), prodiName); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_PRODI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_PRODI, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}

// GetProdi godoc
// @Summary Get Prodi
// @Description Mendapatkan data prodi. Jika query id/name diberikan, return 1 data spesifik; tanpa query, return semua data
// @Description
// @Description  **Akses:** Logged User
// @Tags prodi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string false "Nama Prodi" example(teknik-elektro)
// @Param id query int false "ID Prodi" example(1)
// @Success      200      {object}  utils.Response[[]dto.ProdiResponse,any]
// @Failure      400      {object}  swagger.ErrGetProdiFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      500      {object}  swagger.ErrUpdateProdiInternalServer
// @Router /api/prodi [get]
func (c *prodiController) GetProdi(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var prodiQuery dto.ProdiQuery
	if err := ctx.ShouldBindQuery(&prodiQuery); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PRODI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if prodiQuery.ID != 0 && prodiQuery.Name != "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PRODI, dto.ErrQueryParams.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if prodiQuery.Name != "" {
		prodiQuery.Name = helpers.NormalizeString(prodiQuery.Name)
		prodi, err := c.prodiService.GetProdiByName(ctx.Request.Context(), prodiQuery.Name)
		if err != nil {
			status := http.StatusBadRequest
			if errors.Is(err, constants.ErrInternalErr) {
				status = http.StatusInternalServerError
			}
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PRODI, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(status, res)
			return
		}
		res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_PRODI, prodi, path)
		ctx.JSON(http.StatusOK, res)
		return
	}

	if prodiQuery.ID != 0 {
		prodi, err := c.prodiService.GetProdiById(ctx.Request.Context(), prodiQuery.ID)
		if err != nil {
			status := http.StatusBadRequest
			if errors.Is(err, constants.ErrInternalErr) {
				status = http.StatusInternalServerError
			}
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PRODI, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(status, res)
			return
		}
		res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_PRODI, prodi, path)
		ctx.JSON(http.StatusOK, res)
		return
	}

	prodis, err := c.prodiService.GetAllProdi(ctx.Request.Context())
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PRODI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_PRODI, prodis, path)
	ctx.JSON(http.StatusOK, res)
}

// GetProdiByJurusan godoc
// @Summary Get Prodi by Jurusan
// @Description Mendapatkan daftar prodi berdasarkan nama jurusan
// @Description
// @Description  **Akses:** Logged User
// @Tags prodi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param jurusan_name path string true "Nama Jurusan" example(teknik-elektro)
// @Param page query int false "Halaman (default: 1, per halaman: 10)" example(1)
// @Success      200      {object}  utils.Response[[]dto.ProdiResponse,any]
// @Failure      400      {object}  swagger.ErrGetProdiFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      500      {object}  swagger.ErrGetProdiInternalServer
// @Router /api/prodi/jurusan/{jurusan_name} [get]
func (c *prodiController) GetProdiByJurusan(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var jurusanURI dto.JurusanURI
	if err := ctx.ShouldBindUri(&jurusanURI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PRODI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	prodis, err := c.prodiService.GetProdiByJurusanName(ctx.Request.Context(), jurusanURI.JurusanName)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PRODI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_PRODI, prodis, path)
	ctx.JSON(http.StatusOK, res)
}
