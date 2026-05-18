package controller

import (
	"errors"
	"net/http"
	"web-hosting/internal/modules/kurikulum/dto"
	"web-hosting/internal/modules/kurikulum/service"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/utils"

	_ "web-hosting/internal/package/swagger"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type (
	PivotController interface {
		CreatePivot(ctx *gin.Context)
		UpdatePivot(ctx *gin.Context)
		DeletePivot(ctx *gin.Context)
	}

	pivotController struct {
		kurikulumMkService service.KurikulumMKService
		db                 *gorm.DB
	}
)

func NewPivotController(injector do.Injector, kurikulumMkService service.KurikulumMKService, db *gorm.DB) PivotController {
	return &pivotController{
		kurikulumMkService: kurikulumMkService,
		db:                 db,
	}
}

// CreateKurikulumPivot godoc
// @Summary Create Kurikulum Pivot mata kuliah Baru
// @Description Menambahkan kurikulum baru ke sistem
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'KurikulumName' Error:..."`
// @Description  - `400` kurikulum dengan nama tersebut sudah ada -> `message: "failed to create kurikulum", error: "kurikulum already exists"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` user tidak memiliki akses -> `message: "kurikulum anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to create kurikulum", error: "Internal Error"`
// @Tags kurikulum
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.PivotCreateRequest true "Kurikulum Request"
// @Success      201      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrCreateKurikulumPivotFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrCreateKurikulumPivotInternalServer
// @Router /api/kurikulum/mata-kuliah [post]
func (c *pivotController) CreatePivot(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var req dto.PivotCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_KURIKULUM, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := c.kurikulumMkService.CreateKurikulumMK(ctx.Request.Context(), req); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_KURIKULUM, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_KURIKULUM, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}

// UpdateKurikulumPivot godoc
// @Summary Update Kurikulum Pivot
// @Description Mengupdate kurikulum yang sudah ada
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter URI tidak valid -> `message: "failed to validate kurikulum URI", error: "Key: 'KurikulumName' Error:..."`
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'KurikulumName' Error:..."`
// @Description  - `400` kurikulum dengan nama tersebut tidak ditemukan -> `message: "failed to update kurikulum", error: "kurikulum not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` kurikulum user tidak memiliki akses -> `message: "kurikulum anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to update kurikulum", error: "Internal Error"`
// @Tags kurikulum
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param kode path string true "Kurikulum kode"
// @Param mk_kode path string true "mata kuliah kode"
// @Param request body dto.KurikulumUpdateRequest true "Kurikulum Request"
// @Success      200      {object}  utils.Response[dto.PivotResponse,any]
// @Failure      400      {object}  swagger.ErrUpdateKurikulumPivotFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrUpdateKurikulumPivotInternalServer
// @Router /api/kurikulum/{kode}/mata-kuliah/{mk_kode} [put]
func (c *pivotController) UpdatePivot(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.PivotURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KURIKULUM, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var req dto.PivotUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KURIKULUM, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	data, err := c.kurikulumMkService.UpdateKurikulumMK(ctx.Request.Context(), URI.KurikulumKode, URI.MkKode, req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KURIKULUM, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_KURIKULUM, data, path)
	ctx.JSON(http.StatusOK, res)
}

// DeleteKurikulumPivot godoc
// @Summary Delete Kurikulum Pivot
// @Description delete kurikulum pivot yang sudah ada
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter URI tidak valid -> `message: "failed to validate kurikulum URI", error: "Key: 'KurikulumName' Error:..."`
// @Description  - `400` kurikulum pivot dengan kode tersebut tidak ditemukan -> `message: "failed to Delete kurikulum pivot", error: "kurikulum pivot not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` kurikulum pivot user tidak memiliki akses -> `message: "kurikulum pivot anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to Delete kurikulum pivot", error: "Internal Error"`
// @Tags kurikulum
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param kode path string true "Kurikulum kode"
// @Param mk_kode path string true "mata kuliah kode"
// @Success      200      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrDeleteKurikulumPivotFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrDeleteKurikulumPivotInternalServer
// @Router /api/kurikulum/{kode}/mata-kuliah/{mk_kode} [delete]
func (c *pivotController) DeletePivot(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.PivotURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_KURIKULUM, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.kurikulumMkService.DeleteKurikulumMK(ctx.Request.Context(), URI.KurikulumKode, URI.MkKode)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_KURIKULUM, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_KURIKULUM, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}
