package controller

import (
	"errors"
	"net/http"
	"web-hosting/internal/modules/akademik/dto"
	"web-hosting/internal/modules/akademik/service"
	"web-hosting/internal/modules/akademik/validation"
	"web-hosting/internal/package/constants"
	_ "web-hosting/internal/package/swagger"
	"web-hosting/internal/package/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type (
	TahunAkademikController interface {
		CreateTahunAkademik(ctx *gin.Context)
		UpdateTahunAkademik(ctx *gin.Context)
		DeleteTahunAkademik(ctx *gin.Context)
		GetTahunAkademik(ctx *gin.Context)
		GetAllTahunAkademik(ctx *gin.Context)
		GetTahunAkademikByStatus(ctx *gin.Context)
	}

	tahunAkademikController struct {
		akademikService   service.TahunAkademikService
		akademikValidator *validation.TahunAkademikValidator
		db                *gorm.DB
	}
)

func NewTahunAkademikController(injector do.Injector, akademikService service.TahunAkademikService, db *gorm.DB) TahunAkademikController {
	akademikValidator := validation.NewTahunAkademikValidator()
	return &tahunAkademikController{
		akademikService:   akademikService,
		akademikValidator: akademikValidator,
		db:                db,
	}
}

// CreateTahunAkademik godoc
// @Summary Create Tahun Akademik Baru
// @Description Menambahkan tahun akademik baru ke sistem
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'TahunAkademikName' Error:..."`
// @Description  - `400` tahun akademik dengan nama tersebut sudah ada -> `message: "failed to create tahun akademik", error: "tahun akademik already exists"`
// @Description  - `400` tahun akademik dengan tahun awal > akhir -> `message: "failed to create tahun akademik", error: "invalid tahun awal/akhir akademik "`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` user tidak memiliki akses -> `message: "tahun akademik anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to create tahun akademik", error: "Internal Error"`
// @Tags tahun-akademik
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body swagger.AkademikCreateRequest true "Tahun Akademik Request"
// @Success      201      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrCreateTahunAkademikFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrCreateTahunAkademikInternalServer
// @Router /api/tahun-akademik [post]
func (c *tahunAkademikController) CreateTahunAkademik(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var req dto.AkademikCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_TAHUN_AKADEMIK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.akademikService.CreateTahunAkademik(ctx.Request.Context(), req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_TAHUN_AKADEMIK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_TAHUN_AKADEMIK, any(nil), path)
	ctx.JSON(http.StatusCreated, res)
}

// UpdateTahunAkademik godoc
// @Summary Update Tahun Akademik
// @Description Mengupdate tahun akademik yang sudah ada
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter Query tidak valid -> `message: "failed to validate tahun akademik Query", error: "Key: 'TahunAkademikName' Error:..."`
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'TahunAkademikName' Error:..."`
// @Description  - `400` tahun akademik dengan tahun awal > akhir -> `message: "failed to create tahun akademik", error: "invalid tahun awal/akhir akademik "`
// @Description  - `400` tahun akademik dengan nama tersebut tidak ditemukan -> `message: "failed to update tahun akademik", error: "tahun akademik not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` jurusan user tidak memiliki akses -> `message: "jurusan anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to update jurusan", error: "Internal Error"`
// @Tags tahun-akademik
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Tahun Akademik ID"
// @Param request body swagger.AkademikUpdateRequest true "Tahun Akademik Request"
// @Success      200      {object}  utils.Response[swagger.AkademikResponse,any]
// @Failure      400      {object}  swagger.ErrUpdateTahunAkademikFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrUpdateTahunAkademikInternalServer
// @Router /api/tahun-akademik/{id} [put]
func (c *tahunAkademikController) UpdateTahunAkademik(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var uri dto.AkademikIdUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_TAHUN_AKADEMIK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var req dto.AkademikUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_TAHUN_AKADEMIK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.akademikService.UpdateTahunAkademik(ctx.Request.Context(), uri.ID, req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_TAHUN_AKADEMIK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_TAHUN_AKADEMIK, result, path)
	ctx.JSON(http.StatusOK, res)
}

// DeleteTahunAkademik godoc
// @Summary Delete Tahun Akademik
// @Description delete tahun akademik yang sudah ada
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter Query tidak valid -> `message: "failed to validate tahun akademik Query", error: "Key: 'TahunAkademikName' Error:..."`
// @Description  - `400` tahun akademik dengan nama tersebut tidak ditemukan -> `message: "failed to Delete tahun akademik", error: "tahun akademik not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` jurusan user tidak memiliki akses -> `message: "jurusan anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to Delete tahun akademik", error: "Internal Error"`
// @Tags tahun-akademik
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Tahun Akademik ID"
// @Success      200      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrDeleteTahunAkademikFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrDeleteTahunAkademikInternalServer
// @Router /api/tahun-akademik/{id} [delete]
func (c *tahunAkademikController) DeleteTahunAkademik(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var uri dto.AkademikIdUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_TAHUN_AKADEMIK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.akademikService.DeleteTahunAkademik(ctx.Request.Context(), uri.ID)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_TAHUN_AKADEMIK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_TAHUN_AKADEMIK, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}

// GetTahunAkademikByID godoc
// @Summary get Tahun Akademik
// @Description melihat tahun akademik yang sudah ada
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter Query tidak valid -> `message: "failed to validate tahun akademik Query", error: "Key: 'TahunAkademikName' Error:..."`
// @Description  - `400` tahun akademik dengan nama tersebut tidak ditemukan -> `message: "failed to update tahun akademik", error: "tahun akademik not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` jurusan user tidak memiliki akses -> `message: "jurusan anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to update jurusan", error: "Internal Error"`
// @Tags tahun-akademik
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Tahun Akademik ID"
// @Success      200      {object}  utils.Response[swagger.AkademikResponse,any]
// @Failure      400      {object}  swagger.ErrGetTahunAkademikFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrGetTahunAkademikInternalServer
// @Router /api/tahun-akademik/{id} [get]
func (c *tahunAkademikController) GetTahunAkademik(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var uri dto.AkademikIdUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TAHUN_AKADEMIK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.akademikService.GetTahunAkademikByID(ctx.Request.Context(), uri.ID)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TAHUN_AKADEMIK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_TAHUN_AKADEMIK, result, path)
	ctx.JSON(http.StatusOK, res)
}

// GetAllTahunAkademik godoc
// @Summary get semua Tahun Akademik
// @Description melihat tahun akademik yang sudah ada
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter Query tidak valid -> `message: "failed to validate tahun akademik Query", error: "Key: 'TahunAkademikName' Error:..."`
// @Description  - `400` tahun akademik dengan nama tersebut tidak ditemukan -> `message: "failed to update tahun akademik", error: "tahun akademik not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` jurusan user tidak memiliki akses -> `message: "jurusan anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to update jurusan", error: "Internal Error"`
// @Tags tahun-akademik
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success      200      {object}  utils.Response[[]swagger.AkademikResponse,any]
// @Failure      400      {object}  swagger.ErrGetTahunAkademikFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrGetTahunAkademikInternalServer
// @Router /api/tahun-akademik [get]
func (c *tahunAkademikController) GetAllTahunAkademik(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	result, err := c.akademikService.GetAllTahunAkademik(ctx.Request.Context())
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TAHUN_AKADEMIK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_TAHUN_AKADEMIK, result, path)
	ctx.JSON(http.StatusOK, res)
}

// GetTahunAkademikByStatus godoc
// @Summary get tahun Akademik berdasarkan status
// @Description melihat tahun akademik yang sudah ada
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter Query tidak valid -> `message: "failed to validate tahun akademik Query", error: "Key: 'TahunAkademikName' Error:..."`
// @Description  - `400` tahun akademik dengan nama tersebut tidak ditemukan -> `message: "failed to update tahun akademik", error: "tahun akademik not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` jurusan user tidak memiliki akses -> `message: "jurusan anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to update jurusan", error: "Internal Error"`
// @Tags tahun-akademik
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param status path string true "Status tahun akademik"
// @Success      200      {object}  utils.Response[[]swagger.AkademikResponse,any]
// @Failure      400      {object}  swagger.ErrGetTahunAkademikFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrGetTahunAkademikInternalServer
// @Router /api/tahun-akademik/status/{status} [get]
func (c *tahunAkademikController) GetTahunAkademikByStatus(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var statusQuery dto.AkademikStatusUri
	if err := ctx.ShouldBindUri(&statusQuery); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TAHUN_AKADEMIK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.akademikService.GetTahunAkademikByStatus(ctx.Request.Context(), statusQuery.Status)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TAHUN_AKADEMIK, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_TAHUN_AKADEMIK, result, path)
	ctx.JSON(http.StatusOK, res)
}
