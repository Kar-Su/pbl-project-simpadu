package controller

import (
	"errors"
	"net/http"
	"web-hosting/internal/modules/kelas/dto"
	"web-hosting/internal/modules/kelas/service"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type (
	KelasController interface {
		CreateKelas(ctx *gin.Context)
		UpdateKelas(ctx *gin.Context)
		DeleteKelas(ctx *gin.Context)
		GetKelasByID(ctx *gin.Context)
		GetKelasByProdiName(ctx *gin.Context)
	}

	kelasController struct {
		db           *gorm.DB
		kelasService service.KelasService
	}
)

// CreateKelas godoc
// @Summary Create Kelas Baru
// @Description Menambahkan kelas baru ke sistem
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'KelasName' Error:..."`
// @Description  - `400` kelas dengan nama tersebut sudah ada -> `message: "failed to create kelas", error: "kelas already exists"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` user tidak memiliki akses -> `message: "kelas anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to create kelas", error: "Internal Error"`
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
func NewKelasController(injector do.Injector, db *gorm.DB, kelasService service.KelasService) KelasController {
	return &kelasController{
		db:           db,
		kelasService: kelasService,
	}
}

func (c *kelasController) CreateKelas(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var req dto.KelasCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.kelasService.CreateKelas(ctx.Request.Context(), req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_KELAS, any(nil), path)
	ctx.JSON(http.StatusCreated, res)
}

func (c *kelasController) UpdateKelas(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.KelasIdURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var req dto.KelasUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	data, err := c.kelasService.UpdateKelas(ctx.Request.Context(), URI.KelasID, req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_KELAS, data, path)
	ctx.JSON(http.StatusOK, res)
}

func (c *kelasController) DeleteKelas(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.KelasIdURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.kelasService.DeleteKelas(ctx.Request.Context(), URI.KelasID)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_KELAS, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}

func (c *kelasController) GetKelasByID(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.KelasIdURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	data, err := c.kelasService.GetKelasByID(ctx.Request.Context(), URI.KelasID)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_KELAS, data, path)
	ctx.JSON(http.StatusOK, res)
}

func (c *kelasController) GetKelasByProdiName(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.ProdiNameURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	data, err := c.kelasService.GetKelasByProdiName(ctx.Request.Context(), URI.Name)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_KELAS, data, path)
	ctx.JSON(http.StatusOK, res)
}
