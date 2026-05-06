package controller

import (
	"errors"
	"net/http"
	_ "web-hosting/internal/database/entities"
	"web-hosting/internal/modules/jurusan/dto"
	"web-hosting/internal/modules/jurusan/service"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/helpers"
	_ "web-hosting/internal/package/swagger"
	"web-hosting/internal/package/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type JurusanController interface {
	CreateJurusan(ctx *gin.Context)
	UpdateJurusan(ctx *gin.Context)
	DeleteJurusan(ctx *gin.Context)
	GetJurusan(ctx *gin.Context)
}

type jurusanController struct {
	jurusanService service.JurusanService
	db             *gorm.DB
}

func NewJurusanController(injector do.Injector, jurusanService service.JurusanService, db *gorm.DB) *jurusanController {
	return &jurusanController{
		jurusanService: jurusanService,
		db:             db,
	}
}

// CreateJurusan godoc
// @Summary Create Jurusan Baru
// @Description Menambahkan jurusan baru ke sistem
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'JurusanName' Error:..."`
// @Description  - `400` jurusan dengan nama tersebut sudah ada -> `message: "failed to create jurusan", error: "jurusan already exists"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` user tidak memiliki akses -> `message: "jurusan anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to create jurusan", error: "Internal Error"`
// @Tags jurusan
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.JurusanRequest true "Jurusan Request"
// @Success      201      {object}  utils.Response[entities.Jurusan,any]
// @Failure      400      {object}  swagger.ErrCreateJurusanFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrCreateJurusanInternalServer
// @Router /api/jurusan [post]
func (c *jurusanController) CreateJurusan(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var req dto.JurusanRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_JURUSAN, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	jurusan, err := c.jurusanService.CreateJurusan(ctx.Request.Context(), req)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_JURUSAN, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_JURUSAN, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_JURUSAN, jurusan, path)
	ctx.JSON(http.StatusOK, res)
}

// UpdateJurusan godoc
// @Summary Update Jurusan
// @Description Mengupdate jurusan yang sudah ada
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter Query tidak valid -> `message: "failed to validate jurusan Query", error: "Key: 'JurusanName' Error:..."`
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'JurusanName' Error:..."`
// @Description  - `400` jurusan dengan nama tersebut tidak ditemukan -> `message: "failed to update jurusan", error: "jurusan not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` jurusan user tidak memiliki akses -> `message: "jurusan anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to update jurusan", error: "Internal Error"`
// @Tags jurusan
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query dto.JurusanNameQuery true "Jurusan Name"
// @Param request body dto.JurusanUpdateRequest true "Jurusan Request"
// @Success      200      {object}  utils.Response[entities.Jurusan,any]
// @Failure      400      {object}  swagger.ErrUpdateJurusanFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrUpdateJurusanInternalServer
// @Router /api/jurusan/ [put]
func (c *jurusanController) UpdateJurusan(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var jurusanNameQuery dto.JurusanNameQuery
	if err := ctx.ShouldBindQuery(&jurusanNameQuery); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_JURUSAN, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	jurusanName := helpers.NormalizeString(jurusanNameQuery.JurusanName)

	var req dto.JurusanUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_JURUSAN, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	updatedJurusan, err := c.jurusanService.UpdateJurusan(ctx.Request.Context(), req, jurusanName)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_JURUSAN, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_JURUSAN, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_JURUSAN, updatedJurusan, path)
	ctx.JSON(http.StatusOK, res)
}

// DeleteJurusan godoc
// @Summary Delete Jurusan
// @Description delete jurusan yang sudah ada
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter Query tidak valid -> `message: "failed to validate jurusan Query", error: "Key: 'JurusanName' Error:..."`
// @Description  - `400` jurusan dengan nama tersebut tidak ditemukan -> `message: "failed to Delete jurusan", error: "jurusan not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` jurusan user tidak memiliki akses -> `message: "jurusan anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to Delete jurusan", error: "Internal Error"`
// @Tags jurusan
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query dto.JurusanNameQuery true "Jurusan Name"
// @Param request body dto.JurusanRequest true "Jurusan Request"
// @Success      200      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrDeleteJurusanFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrDeleteJurusanInternalServer
// @Router /api/jurusan/ [delete]
func (c *jurusanController) DeleteJurusan(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var jurusanNameQuery dto.JurusanNameQuery
	if err := ctx.ShouldBindQuery(&jurusanNameQuery); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_JURUSAN, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	jurusanName := helpers.NormalizeString(jurusanNameQuery.JurusanName)

	if err := c.jurusanService.DeleteJurusan(ctx.Request.Context(), jurusanName); err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_JURUSAN, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_JURUSAN, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_JURUSAN, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}

// GetJurusan godoc
// @Summary get Jurusan
// @Description melihat jurusan yang sudah ada
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter Query tidak valid -> `message: "failed to validate jurusan Query", error: "Key: 'JurusanName' Error:..."`
// @Description  - `400` jurusan dengan nama tersebut tidak ditemukan -> `message: "failed to update jurusan", error: "jurusan not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` jurusan user tidak memiliki akses -> `message: "jurusan anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to update jurusan", error: "Internal Error"`
// @Tags jurusan
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query dto.JurusanQuery false "Jurusan Name Or ID"
// @Success      200      {object}  utils.Response[entities.Jurusan,any]
// @Success      200      {object}  utils.Response[[]entities.Jurusan,any]
// @Failure      400      {object}  swagger.ErrUpdateJurusanFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrUpdateJurusanInternalServer
// @Router /api/jurusan/ [get]
func (c *jurusanController) GetJurusan(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var query dto.JurusanQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_JURUSAN, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var (
		jurusan any
		err     error
		message string
	)

	if query.JurusanID != 0 {
		jurusan, err = c.jurusanService.GetJurusanById(ctx.Request.Context(), query.JurusanID)
		message = dto.MESSAGE_SUCCESS_GET_JURUSAN
	} else if query.JurusanName != "" {
		normalizedName := helpers.NormalizeString(query.JurusanName)
		jurusan, err = c.jurusanService.GetJurusanByName(ctx.Request.Context(), normalizedName)
		message = dto.MESSAGE_SUCCESS_GET_JURUSAN
	} else {
		jurusan, err = c.jurusanService.GetAllJurusan(ctx.Request.Context())
		message = dto.MESSAGE_SUCCESS_GET_ALL
	}

	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_JURUSAN, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(message, jurusan, path)
	ctx.JSON(http.StatusOK, res)
}
