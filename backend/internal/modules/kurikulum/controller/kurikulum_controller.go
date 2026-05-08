package controller

import (
	"errors"
	"net/http"
	"web-hosting/internal/modules/kurikulum/dto"
	"web-hosting/internal/modules/kurikulum/service"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/helpers"
	"web-hosting/internal/package/utils"

	_ "web-hosting/internal/package/swagger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type (
	KurikulumController interface {
		CreateKurikulum(ctx *gin.Context)
		UpdateKurikulum(ctx *gin.Context)
		DeleteKurikulum(ctx *gin.Context)
		GetKurikulum(ctx *gin.Context)
	}
	kurikulumController struct {
		kurikulumService service.KurikulumService
		db               *gorm.DB
	}
)

func NewKurikulumController(injector do.Injector, kurikulumService service.KurikulumService, db *gorm.DB) KurikulumController {
	return &kurikulumController{
		kurikulumService: kurikulumService,
		db:               db,
	}
}

// CreateKurikulum godoc
// @Summary Create Kurikulum Baru
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
// @Param request body dto.KurikulumCreateRequest true "Kurikulum Request"
// @Success      201      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrCreateKurikulumFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrCreateKurikulumInternalServer
// @Router /api/kurikulum [post]
func (c *kurikulumController) CreateKurikulum(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var req dto.KurikulumCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_KURIKULUM, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.kurikulumService.CreateKurikulum(ctx.Request.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_KURIKULUM, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_KURIKULUM, any(nil), path)
	ctx.JSON(http.StatusCreated, res)
}

// UpdateKurikulum godoc
// @Summary Update Kurikulum
// @Description Mengupdate kurikulum yang sudah ada
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter Query tidak valid -> `message: "failed to validate kurikulum Query", error: "Key: 'KurikulumName' Error:..."`
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'KurikulumName' Error:..."`
// @Description  - `400` kurikulum dengan nama tersebut tidak ditemukan -> `message: "failed to update kurikulum", error: "kurikulum not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` role user tidak memiliki akses -> `message: "role anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to update role", error: "Internal Error"`
// @Tags kurikulum
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query dto.KurikulumQuery true "Kurikulum Name"
// @Param request body dto.KurikulumUpdateRequest true "Kurikulum Request"
// @Success      200      {object}  utils.Response[dto.KurikulumResponse,any]
// @Failure      400      {object}  swagger.ErrUpdateKurikulumFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrUpdateKurikulumInternalServer
// @Router /api/kurikulum/ [put]
func (c *kurikulumController) UpdateKurikulum(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var query dto.KurikulumQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KURIKULUM, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	if query.ID != "" && query.Kode != "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KURIKULUM, dto.ErrQueryParams, nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var req dto.KurikulumUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KURIKULUM, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var (
		err    error
		status int
		data   dto.KurikulumResponse
	)

	if query.ID != "" {
		id, err := uuid.Parse(query.ID)
		if err != nil {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KURIKULUM, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}
		data, err = c.kurikulumService.UpdateKurikulumById(ctx.Request.Context(), req, id)

	} else if query.Kode != "" {
		query.Kode = helpers.NormalizeString(query.Kode)
		data, err = c.kurikulumService.UpdateKurikulumByKode(ctx.Request.Context(), req, query.Kode)
	} else {
		err = dto.ErrQueryParams
	}

	if err != nil {
		status = http.StatusBadRequest
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

// DeleteKurikulum godoc
// @Summary Delete Kurikulum
// @Description delete kurikulum yang sudah ada
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter Query tidak valid -> `message: "failed to validate kurikulum Query", error: "Key: 'KurikulumName' Error:..."`
// @Description  - `400` kurikulum dengan nama tersebut tidak ditemukan -> `message: "failed to Delete kurikulum", error: "kurikulum not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` kurikulum user tidak memiliki akses -> `message: "kurikulum anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to Delete kurikulum", error: "Internal Error"`
// @Tags kurikulum
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query dto.KurikulumQuery true "Kurikulum Name"
// @Success      200      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrDeleteKurikulumFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrDeleteKurikulumInternalServer
// @Router /api/kurikulum/ [delete]
func (c *kurikulumController) DeleteKurikulum(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var query dto.KurikulumQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_KURIKULUM, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if query.ID != "" && query.Kode != "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KURIKULUM, dto.ErrQueryParams, nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var (
		err    error
		status int
	)

	if query.ID != "" {
		id, err := uuid.Parse(query.ID)
		if err != nil {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_KURIKULUM, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}
		err = c.kurikulumService.DeleteKurikulumById(ctx.Request.Context(), id)
	} else if query.Kode != "" {
		err = c.kurikulumService.DeleteKurikulumByKode(ctx.Request.Context(), query.Kode)
	} else {
		err = dto.ErrQueryParams
	}

	if err != nil {
		status = http.StatusBadRequest
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

// GetKurikulum godoc
// @Summary get Kurikulum
// @Description melihat kurikulum yang sudah ada
// @Description Pilih salah satu query id/name
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter Query tidak valid -> `message: "failed to validate kurikulum Query", error: "Key: 'KurikulumName' Error:..."`
// @Description  - `400` kurikulum dengan nama tersebut tidak ditemukan -> `message: "failed to Get kurikulum", error: "kurikulum not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` kurikulum user tidak memiliki akses -> `message: "kurikulum anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to Get kurikulum", error: "Internal Error"`
// @Tags kurikulum
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query dto.KurikulumQuery false "Kurikulum Name Or ID (Pilih salah satu)"
// @Success      200      {object}  utils.Response[dto.KurikulumResponse,any]
// @Success      200      {object}  utils.Response[[]dto.KurikulumResponse,any]
// @Failure      400      {object}  swagger.ErrGetKurikulumFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrGetKurikulumInternalServer
// @Router /api/kurikulum/ [get]
func (c *kurikulumController) GetKurikulum(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var query dto.KurikulumQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_KURIKULUM, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var (
		err    error
		status int
		data   any
	)

	if query.ID != "" {
		id, err := uuid.Parse(query.ID)
		if err != nil {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_KURIKULUM, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}
		data, err = c.kurikulumService.GetKurikulumById(ctx.Request.Context(), id)
	} else if query.Kode != "" {
		query.Kode = helpers.NormalizeString(query.Kode)
		data, err = c.kurikulumService.GetKurikulumByKode(ctx.Request.Context(), query.Kode)
	} else {
		data, err = c.kurikulumService.GetAllKurikulum(ctx.Request.Context())
	}

	if err != nil {
		status = http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_KURIKULUM, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_KURIKULUM, data, path)
	ctx.JSON(http.StatusOK, res)
}
