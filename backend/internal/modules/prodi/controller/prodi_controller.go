package controller

import (
	"errors"
	"net/http"
	"web-hosting/internal/modules/prodi/dto"
	"web-hosting/internal/modules/prodi/service"
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
		prodiService service.ProdiService
		db           *gorm.DB
	}
)

func NewProdiController(injector do.Injector, prodiService service.ProdiService, db *gorm.DB) ProdiController {
	return &prodiController{
		prodiService: prodiService,
		db:           db,
	}
}

// CreateProdi godoc
// @Summary Create Prodi Baru
// @Description Menambahkan prodi baru ke sistem
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'ProdiName' Error:..."`
// @Description  - `400` prodi dengan nama tersebut sudah ada -> `message: "failed to create prodi", error: "prodi already exists"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` user tidak memiliki akses -> `message: "prodi anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to create prodi", error: "Internal Error"`
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
// @Description Mengupdate prodi yang sudah ada
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter Query tidak valid -> `message: "failed to validate prodi Query", error: "Key: 'ProdiName' Error:..."`
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'ProdiName' Error:..."`
// @Description  - `400` prodi dengan nama tersebut tidak ditemukan -> `message: "failed to update prodi", error: "prodi not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` prodi user tidak memiliki akses -> `message: "prodi anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to update prodi", error: "Internal Error"`
// @Tags prodi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query dto.ProdiNameQuery true "Prodi Name"
// @Param request body dto.ProdiUpdateRequest true "Prodi Request"
// @Success      200      {object}  utils.Response[dto.ProdiResponse,any]
// @Failure      400      {object}  swagger.ErrUpdateProdiFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrUpdateProdiInternalServer
// @Router /api/prodi/ [put]
func (c *prodiController) UpdateProdi(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var prodiNameQuery dto.ProdiNameQuery
	if err := ctx.ShouldBindQuery(&prodiNameQuery); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_PRODI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	prodiNameQuery.Name = helpers.NormalizeString(prodiNameQuery.Name)

	var req dto.ProdiUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_PRODI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	updatedResponse, err := c.prodiService.UpdateProdi(ctx.Request.Context(), prodiNameQuery.Name, req)
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
// @Description delete prodi yang sudah ada
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter Query tidak valid -> `message: "failed to validate prodi Query", error: "Key: 'ProdiName' Error:..."`
// @Description  - `400` prodi dengan nama tersebut tidak ditemukan -> `message: "failed to Delete prodi", error: "prodi not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` prodi user tidak memiliki akses -> `message: "prodi anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to Delete prodi", error: "Internal Error"`
// @Tags prodi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query dto.ProdiNameQuery true "Prodi Name"
// @Success      200      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrDeleteProdiFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrDeleteProdiInternalServer
// @Router /api/prodi/ [delete]
func (c *prodiController) DeleteProdi(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var prodiNameQuery dto.ProdiNameQuery
	if err := ctx.ShouldBindQuery(&prodiNameQuery); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_PRODI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	prodiNameQuery.Name = helpers.NormalizeString(prodiNameQuery.Name)

	if err := c.prodiService.DeleteProdi(ctx.Request.Context(), prodiNameQuery.Name); err != nil {
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
// @Summary get data prodi
// @Description melihat data prodi yang sudah ada
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter Query tidak valid -> `message: "failed to validate prodi Query", error: "Key: 'ProdiName' Error:..."`
// @Description  - `400` prodi dengan nama tersebut tidak ditemukan -> `message: "failed to update prodi", error: "prodi not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` prodi user tidak memiliki akses -> `message: "prodi anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to update prodi", error: "Internal Error"`
// @Tags prodi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query dto.ProdiQuery false "Prodi Name Or ID"
// @Success      200      {object}  utils.Response[dto.ProdiResponse,any]
// @Success      200      {object}  utils.Response[[]dto.ProdiResponse,any]
// @Failure      400      {object}  swagger.ErrGetProdiFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrUpdateProdiInternalServer
// @Router /api/prodi/ [get]
func (c *prodiController) GetProdi(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var prodiQuery dto.ProdiQuery
	if err := ctx.ShouldBindQuery(&prodiQuery); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PRODI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var (
		prodi   any
		err     error
		message string
	)

	if prodiQuery.Name != "" {
		prodiQuery.Name = helpers.NormalizeString(prodiQuery.Name)
		prodi, err = c.prodiService.GetProdiByName(ctx.Request.Context(), prodiQuery.Name)
		message = dto.MESSAGE_SUCCESS_GET_PRODI
	} else if prodiQuery.ID != 0 {
		prodi, err = c.prodiService.GetProdiById(ctx.Request.Context(), prodiQuery.ID)
		message = dto.MESSAGE_SUCCESS_GET_PRODI
	} else {
		prodi, err = c.prodiService.GetAllProdi(ctx.Request.Context())
		message = dto.MESSAGE_SUCCESS_GET_PRODI
	}
	if err != nil {
		status := http.StatusBadRequest
		message = dto.MESSAGE_FAILED_GET_PRODI
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(message, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(message, prodi, path)
	ctx.JSON(http.StatusOK, res)
}

// GetProdiByJurusan godoc
// @Summary get data prodi berdasarkan nama jurusan
// @Description melihat data prodi yang sudah ada
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter Uri tidak valid -> `message: "failed to validate prodi Uri", error: "Key: 'JurusanName' Error:..."`
// @Description  - `400` jurusan dengan nama tersebut tidak ditemukan -> `message: "failed to get jurusan", error: "jurusan not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` prodi user tidak memiliki akses -> `message: "prodi anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to get prodi", error: "Internal Error"`
// @Tags prodi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param jurusan_name path string true "Jurusan Name" example(teknik-elektro)
// @Success      200      {object}  utils.Response[[]dto.ProdiResponse,any]
// @Failure      400      {object}  swagger.ErrGetProdiFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
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
