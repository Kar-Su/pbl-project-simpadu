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
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'MataKuliahName' Error:..."`
// @Description  - `400` mata kuliah dengan nama tersebut sudah ada -> `message: "failed to create mata-kuliah", error: "mata-kuliah already exists"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` user tidak memiliki akses -> `message: "mata-kuliah anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to create mata-kuliah", error: "Internal Error"`
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
	}

	if err := c.mkService.CreateMk(ctx, req); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_MK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_MK, any(nil), path)
	ctx.JSON(http.StatusCreated, res)
}

// UpdateMataKuliah godoc
// @Summary Update Mata Kuliah
// @Description Mengupdate mata kuliah yang sudah ada
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter Query tidak valid -> `message: "failed to validate mata-kuliah Query", error: "Key: 'MataKuliahName' Error:..."`
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'MataKuliahName' Error:..."`
// @Description  - `400` mata kuliah dengan nama tersebut tidak ditemukan -> `message: "failed to update mata-kuliah", error: "mata-kuliah not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` mata kuliah user tidak memiliki akses -> `message: "mata-kuliah anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to update mata-kuliah", error: "Internal Error"`
// @Tags mata-kuliah
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query dto.MkQuery true "Mata Kuliah Name"
// @Param request body dto.MkUpdateRequest true "Mata Kuliah Request"
// @Success      200      {object}  utils.Response[dto.MkResponse,any]
// @Failure      400      {object}  swagger.ErrUpdateMkFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrUpdateMkInternalServer
// @Router /api/mata-kuliah/ [put]
func (c *mkController) UpdateMk(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var QueryParams dto.MkQuery
	if err := ctx.ShouldBindQuery(&QueryParams); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_MK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	var req dto.MkUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_MK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	var (
		mkResponse dto.MkResponse
		err        error
		status     int
		message    string
	)

	if QueryParams.ID != "" {
		mkId := uuid.MustParse(QueryParams.ID)
		mkResponse, err = c.mkService.UpdateMkById(ctx, mkId, req)
	} else if QueryParams.Kode != "" {
		QueryParams.Kode = helpers.NormalizeString(QueryParams.Kode)
		mkResponse, err = c.mkService.UpdateMkByKode(ctx, QueryParams.Kode, req)
	} else {
		err = dto.ErrQueryParams
	}

	if err != nil {
		status = http.StatusBadRequest
		message = dto.MESSAGE_FAILED_GET_MK
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(message, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	status = http.StatusOK
	message = dto.MESSAGE_SUCCESS_UPDATE_MK
	res := utils.BuildResponseSuccess(message, mkResponse, path)
	ctx.JSON(status, res)
}

// DeleteMataKuliah godoc
// @Summary Delete Mata Kuliah
// @Description delete mata kuliah yang sudah ada
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter Query tidak valid -> `message: "failed to validate mata-kuliah Query", error: "Key: 'MataKuliahName' Error:..."`
// @Description  - `400` mata kuliah dengan nama tersebut tidak ditemukan -> `message: "failed to Delete mata-kuliah", error: "mata-kuliah not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` mata kuliah user tidak memiliki akses -> `message: "mata-kuliah anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to Delete mata-kuliah", error: "Internal Error"`
// @Tags mata-kuliah
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query dto.MkQuery true "Mata Kuliah Name"
// @Success      200      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrDeleteMkFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrDeleteMkInternalServer
// @Router /api/mata-kuliah/ [delete]
func (c *mkController) DeleteMk(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var QueryParams dto.MkQuery
	if err := ctx.ShouldBindQuery(&QueryParams); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_MK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var (
		err     error
		status  int
		message string
	)

	if QueryParams.ID != "" {
		mkId := uuid.MustParse(QueryParams.ID)
		err = c.mkService.DeleteMkById(ctx, mkId)
	} else if QueryParams.Kode != "" {
		QueryParams.Kode = helpers.NormalizeString(QueryParams.Kode)
		err = c.mkService.DeleteMkByKode(ctx, QueryParams.Kode)
	} else {
		err = dto.ErrQueryParams
	}

	if err != nil {
		status = http.StatusBadRequest
		message = dto.MESSAGE_FAILED_DELETE_MK
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(message, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	status = http.StatusOK
	message = dto.MESSAGE_SUCCESS_DELETE_MK
	res := utils.BuildResponseSuccess(message, any(nil), path)
	ctx.JSON(status, res)
}

// GetMataKuliah godoc
// @Summary get mata kuliah
// @Description melihat mata kuliah yang sudah ada
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter Query tidak valid -> `message: "failed to validate mata kuliah Query", error: "Key: 'mkName' Error:..."`
// @Description  - `400` mata kuliah dengan nama tersebut tidak ditemukan -> `message: "failed to update mata kuliah", error: "mata kuliah not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` mata kuliah user tidak memiliki akses -> `message: "mata kuliah anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to update mata kuliah", error: "Internal Error"`
// @Tags mata-kuliah
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param query query dto.MkQuery false "Mata Kuliah Name Or ID (Pilih salah satu)"
// @Success      200      {object}  utils.Response[dto.MkResponse,any]
// @Success      200      {object}  utils.Response[[]dto.MkResponse,any]
// @Failure      400      {object}  swagger.ErrUpdateMkFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrUpdateMkInternalServer
// @Router /api/mata-kuliah/ [get]
func (c *mkController) GetMk(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var QueryParams dto.MkQuery
	if err := ctx.ShouldBindQuery(&QueryParams); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_MK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var (
		err     error
		status  int
		message string
		result  any
	)

	if QueryParams.ID != "" && QueryParams.Kode != "" {
		err = dto.ErrQueryParams
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_MK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if QueryParams.ID != "" {
		mkId := uuid.MustParse(QueryParams.ID)
		result, err = c.mkService.GetMkById(ctx, mkId)
	} else if QueryParams.Kode != "" {
		QueryParams.Kode = helpers.NormalizeString(QueryParams.Kode)
		result, err = c.mkService.GetMkByKode(ctx, QueryParams.Kode)
	} else {
		result, err = c.mkService.GetAllMk(ctx)
	}

	if err != nil {
		status = http.StatusInternalServerError
		message = dto.MESSAGE_FAILED_GET_MK
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(message, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	status = http.StatusOK
	message = dto.MESSAGE_SUCCESS_GET_MK
	res := utils.BuildResponseSuccess(message, result, path)
	ctx.JSON(status, res)
}
