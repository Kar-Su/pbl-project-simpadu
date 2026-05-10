package controller

import (
	"errors"
	"net/http"
	"web-hosting/internal/modules/kelas/dto"
	"web-hosting/internal/modules/kelas/service"
	"web-hosting/internal/package/constants"
	_ "web-hosting/internal/package/swagger"
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

// AssignMahasiswaToKelas godoc
// @Summary Assign Mahasiswa to Kelas
// @Description Assign mahasiswa ke kelas
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'KelasName' Error:..."`
// @Description  - `400` kelas dengan id tersebut sudah ada -> `message: "failed to create kelas", error: "kelas already exists"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` user tidak memiliki akses -> `message: "kelas anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to create kelas", error: "Internal Error"`
// @Tags kelas
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.KelasMahasiswaCreateRequest true "Kelas Request"
// @Success      201      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrCreateKelasFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrCreateKelasInternalServer
// @Router /api/kelas/mahasiswa [post]
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

// RemoveMahasiswaFromKelas godoc
// @Summary Remove mahasiswa dari kelas
// @Description delete kelas yang sudah ada
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter tidak valid -> `message: "failed to validate parameter", error: "Key: 'param' Error:..."`
// @Description  - `400` kelas dengan id tersebut tidak ditemukan -> `message: "failed to Delete kelas", error: "kelas not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` kelas user tidak memiliki akses -> `message: "kelas anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to Delete kelas", error: "Internal Error"`
// @Tags kelas
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param kelas_id path string true "kelas ID"
// @Param mahasiswa_id path string true "mahasiswa ID"
// @Success      200      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrDeleteKelasFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrDeleteKelasInternalServer
// @Router /api/kelas/{kelas_id}/mahasiswa/{mahasiswa_id} [delete]
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

// GetKelasByMahasiswa godoc
// @Summary get Kelas
// @Description melihat mahasiswa terdaftar di kelas mana saja
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter tidak valid -> `message: "failed to validate parameter", error: "Key: 'param' Error:..."`
// @Description  - `400` kelas dengan id tersebut tidak ditemukan -> `message: "failed to Get kelas", error: "kelas not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` kelas user tidak memiliki akses -> `message: "kelas anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to Get kelas", error: "Internal Error"`
// @Tags kelas
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param mahasiswa_id path string true "mahasiswa ID"
// @Success      200      {object}  utils.Response[[]dto.KelasMahasiswaResponse,any]
// @Failure      400      {object}  swagger.ErrGetKelasFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrGetKelasInternalServer
// @Router /api/kelas/mahasiswa/{mahasiswa_id} [get]
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

// GetMahasiswaByKelas godoc
// @Summary get Kelas
// @Description melihat mahasiswa yang terdaftar dalam kelas
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter tidak valid -> `message: "failed to validate parameter", error: "Key: 'param' Error:..."`
// @Description  - `400` kelas dengan id tersebut tidak ditemukan -> `message: "failed to Get kelas", error: "kelas not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` kelas user tidak memiliki akses -> `message: "kelas anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to Get kelas", error: "Internal Error"`
// @Tags kelas
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param kelas_id path string true "kelas ID"
// @Success      200      {object}  utils.Response[[]dto.KelasMahasiswaResponse,any]
// @Failure      400      {object}  swagger.ErrGetKelasFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrGetKelasInternalServer
// @Router /api/kelas/{kelas_id}/mahasiswa [get]
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
