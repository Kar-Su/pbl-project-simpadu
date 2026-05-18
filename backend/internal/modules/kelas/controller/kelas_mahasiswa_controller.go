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
	"github.com/google/uuid"
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
// @Description Assign mahasiswa ke kelas berdasarkan kelas_id di path
// @Description
// @Description  **Akses:** Admin Akademik
// @Tags kelas
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param kelas_id path string true "Kelas ID (UUID)"
// @Param request body dto.KelasMahasiswaCreateRequest true "Kelas Mahasiswa Request"
// @Success      201      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrCreateKelasFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrCreateKelasInternalServer
// @Router /api/kelas/{kelas_id}/mahasiswa [post]
func (c *kelasMahasiswaController) AssignMahasiswaToKelas(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var kelasURI dto.KelasIdURI
	if err := ctx.ShouldBindUri(&kelasURI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_ASSIGN_MAHASISWA_TO_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	kelasID, err := uuid.Parse(kelasURI.KelasID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_ASSIGN_MAHASISWA_TO_KELAS, "invalid kelas id", nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var req dto.KelasMahasiswaCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_ASSIGN_MAHASISWA_TO_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// Override kelas_id dari path param (lebih RESTful)
	req.KelasID = kelasID

	if err := c.pivotService.Create(ctx.Request.Context(), req); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_ASSIGN_MAHASISWA_TO_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_ASSIGN_MAHASISWA_TO_KELAS, any(nil), path)
	ctx.JSON(http.StatusCreated, res)
}

// RemoveMahasiswaFromKelas godoc
// @Summary Remove mahasiswa dari kelas
// @Description Menghapus mahasiswa dari kelas berdasarkan kelas_id dan mahasiswa_id
// @Description
// @Description  **Akses:** Admin Akademik
// @Tags kelas
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param kelas_id path string true "Kelas ID (UUID)"
// @Param mahasiswa_id path string true "Mahasiswa ID (UUID)"
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
		return
	}

	kelasID, err := uuid.Parse(URI.KelasID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REMOVE_MAHASISWA_FROM_KELAS, "invalid kelas id", nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	mahasiswaID, err := uuid.Parse(URI.MahasiswaID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REMOVE_MAHASISWA_FROM_KELAS, "invalid mahasiswa id", nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := c.pivotService.Delete(ctx.Request.Context(), mahasiswaID, kelasID); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REMOVE_MAHASISWA_FROM_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REMOVE_MAHASISWA_FROM_KELAS, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}

// GetAllKelasMahasiswa godoc
// @Summary Get semua kelas untuk mahasiswa
// @Description Mendapatkan semua kelas yang diikuti oleh mahasiswa berdasarkan mahasiswa_id
// @Description
// @Description  **Akses:** Logged User
// @Tags kelas
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param mahasiswa_id path string true "Mahasiswa ID (UUID)"
// @Success      200      {object}  utils.Response[[]dto.KelasMahasiswaResponse,any]
// @Failure      400      {object}  swagger.ErrGetKelasFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      500      {object}  swagger.ErrGetKelasInternalServer
// @Router /api/mahasiswa/{mahasiswa_id}/kelas [get]
func (c *kelasMahasiswaController) GetAllKelasMahasiswa(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.MahasiswaIdURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_PIVOT, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	mahasiswaID, err := uuid.Parse(URI.MahasiswaID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_PIVOT, "invalid mahasiswa id", nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	data, err := c.pivotService.GetAllKelasMahasiswa(ctx.Request.Context(), mahasiswaID)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_PIVOT, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DATA_PIVOT, data, path)
	ctx.JSON(http.StatusOK, res)
}

// GetMahasiswaByKelas godoc
// @Summary Get mahasiswa dalam kelas
// @Description Mendapatkan daftar mahasiswa yang terdaftar dalam kelas
// @Description
// @Description  **Akses:** Logged User
// @Tags kelas
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param kelas_id path string true "Kelas ID (UUID)"
// @Success      200      {object}  utils.Response[[]dto.KelasMahasiswaResponse,any]
// @Failure      400      {object}  swagger.ErrGetKelasFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      500      {object}  swagger.ErrGetKelasInternalServer
// @Router /api/kelas/{kelas_id}/mahasiswa [get]
func (c *kelasMahasiswaController) GetMahasiswaByKelas(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.KelasIdURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_PIVOT, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	kelasID, err := uuid.Parse(URI.KelasID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_PIVOT, "invalid kelas id", nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	data, err := c.pivotService.GetMahasiswaByKelasId(ctx.Request.Context(), kelasID)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_PIVOT, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DATA_PIVOT, data, path)
	ctx.JSON(http.StatusOK, res)
}
