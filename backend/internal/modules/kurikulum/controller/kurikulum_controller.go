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
		status := http.StatusBadRequest
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
// @Summary Update Kurikulum berdasarkan kode
// @Description Mengupdate kurikulum berdasarkan kode di path parameter
// @Description
// @Description  **Akses:** Admin Akademik
// @Tags kurikulum
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param kode path string true "Kode Kurikulum" example(myhutao-2024)
// @Param request body dto.KurikulumUpdateRequest true "Kurikulum Request"
// @Success      200      {object}  utils.Response[dto.KurikulumResponse,any]
// @Failure      400      {object}  swagger.ErrUpdateKurikulumFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrUpdateKurikulumInternalServer
// @Router /api/kurikulum/{kode} [put]
func (c *kurikulumController) UpdateKurikulum(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var uriParam dto.KurikulumKodeURI
	if err := ctx.ShouldBindUri(&uriParam); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KURIKULUM, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	kode := helpers.NormalizeString(uriParam.Kode)

	var req dto.KurikulumUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KURIKULUM, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	data, err := c.kurikulumService.UpdateKurikulumByKode(ctx.Request.Context(), req, kode)
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

// DeleteKurikulum godoc
// @Summary Delete Kurikulum berdasarkan kode
// @Description Menghapus kurikulum berdasarkan kode di path parameter
// @Description
// @Description  **Akses:** Admin Akademik
// @Tags kurikulum
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param kode path string true "Kode Kurikulum" example(myhutao-2024)
// @Success      200      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrDeleteKurikulumFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrDeleteKurikulumInternalServer
// @Router /api/kurikulum/{kode} [delete]
func (c *kurikulumController) DeleteKurikulum(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var uriParam dto.KurikulumKodeURI
	if err := ctx.ShouldBindUri(&uriParam); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_KURIKULUM, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	kode := helpers.NormalizeString(uriParam.Kode)

	if err := c.kurikulumService.DeleteKurikulumByKode(ctx.Request.Context(), kode); err != nil {
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

// GetKurikulum godoc
// @Summary Get Kurikulum
// @Description Mendapatkan data kurikulum. Jika path kode diberikan, return 1 data; tanpa kode, return semua dengan pagination
// @Description
// @Description  **Akses:** Logged User
// @Tags kurikulum
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param kode path string false "Kode Kurikulum (opsional)" example(myhutao-2024)
// @Param page query int false "Halaman (default: 1, per halaman: 10)" example(1)
// @Success      200      {object}  utils.Response[dto.KurikulumResponse,any]
// @Success      200      {object}  utils.Response[utils.PaginatedData[[]dto.KurikulumResponse],any]
// @Failure      400      {object}  swagger.ErrGetKurikulumFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      500      {object}  swagger.ErrGetKurikulumInternalServer
// @Router /api/kurikulum [get]
// @Router /api/kurikulum/{kode} [get]
func (c *kurikulumController) GetKurikulum(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	// Cek apakah ada path param kode
	kodeParam := ctx.Param("kode")
	if kodeParam != "" && kodeParam != "/" {
		kode := helpers.NormalizeString(kodeParam)
		data, err := c.kurikulumService.GetKurikulumByKode(ctx.Request.Context(), kode)
		if err != nil {
			status := http.StatusBadRequest
			if errors.Is(err, constants.ErrInternalErr) {
				status = http.StatusInternalServerError
			}
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_KURIKULUM, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(status, res)
			return
		}
		res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_KURIKULUM, data, path)
		ctx.JSON(http.StatusOK, res)
		return
	}

	// Pagination untuk get all
	var pageQuery utils.PaginationQuery
	if err := ctx.ShouldBindQuery(&pageQuery); err != nil || pageQuery.Page <= 0 {
		pageQuery.Page = 1
	}

	kurikulums, total, err := c.kurikulumService.GetAllKurikulumPaginated(ctx.Request.Context(), pageQuery.Page)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_KURIKULUM, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	paginated := utils.BuildPaginatedResponse(kurikulums, pageQuery.Page, total)
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_KURIKULUM, paginated, path)
	ctx.JSON(http.StatusOK, res)
}
