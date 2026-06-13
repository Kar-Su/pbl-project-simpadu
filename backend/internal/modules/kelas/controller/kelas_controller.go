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
	KelasController interface {
		CreateKelas(ctx *gin.Context)
		UpdateKelas(ctx *gin.Context)
		DeleteKelas(ctx *gin.Context)
		GetKelasByID(ctx *gin.Context)
		GetKelasByProdiName(ctx *gin.Context)
		GetAllKelas(ctx *gin.Context)
	}

	kelasController struct {
		db           *gorm.DB
		kelasService service.KelasService
	}
)

func NewKelasController(injector do.Injector, db *gorm.DB, kelasService service.KelasService) KelasController {
	return &kelasController{
		db:           db,
		kelasService: kelasService,
	}
}

// CreateKelas godoc
// @Summary Create Kelas Baru
// @Description Menambahkan kelas baru ke sistem
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
// @Param request body dto.KelasCreateRequest true "Kelas Request"
// @Success      201      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrCreateKelasFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrCreateKelasInternalServer
// @Router /api/kelas [post]
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

// UpdateKelas godoc
// @Summary Update kelas
// @Description Mengupdate kelas yang sudah ada
// @Description
// @Description  **Akses:** Admin Akademik
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter tidak valid -> `message: "failed to validate parameter", error: "Key: 'param' Error:..."`
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'Kelas' Error:..."`
// @Description  - `400` kelas dengan id tersebut tidak ditemukan -> `message: "failed to update kelas", error: "kelas not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` role user tidak memiliki akses -> `message: "role anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to update role", error: "Internal Error"`
// @Tags kelas
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param kelas_id path string true "kelas ID"
// @Param request body dto.KelasUpdateRequest true "kelas Request"
// @Success      200      {object}  utils.Response[dto.KelasResponse,any]
// @Failure      400      {object}  swagger.ErrUpdateKelasFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrUpdateKelasInternalServer
// @Router /api/kelas/{kelas_id} [put]
func (c *kelasController) UpdateKelas(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.KelasIdURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	kelasID, err := uuid.Parse(URI.KelasID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KELAS, "invalid kelas id", nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var req dto.KelasUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	data, err := c.kelasService.UpdateKelas(ctx.Request.Context(), kelasID, req)
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

// DeleteKelas godoc
// @Summary Delete kelas
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
// @Success      200      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrDeleteKelasFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrDeleteKelasInternalServer
// @Router /api/kelas/{kelas_id} [delete]
func (c *kelasController) DeleteKelas(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.KelasIdURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	kelasID, err := uuid.Parse(URI.KelasID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_KELAS, "invalid kelas id", nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := c.kelasService.DeleteKelas(ctx.Request.Context(), kelasID); err != nil {
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

// GetKelasByID godoc
// @Summary get Kelas
// @Description melihat kelas yang sudah ada
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
// @Success      200      {object}  utils.Response[dto.KelasResponse,any]
// @Failure      400      {object}  swagger.ErrGetKelasFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrGetKelasInternalServer
// @Router /api/kelas/{kelas_id} [get]
func (c *kelasController) GetKelasByID(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.KelasIdURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	kelasID, err := uuid.Parse(URI.KelasID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_KELAS, "invalid kelas id", nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	data, err := c.kelasService.GetKelasByID(ctx.Request.Context(), kelasID)
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

// GetKelasByProdiName godoc
// @Summary Get Kelas by Prodi (Paginated)
// @Description Mendapatkan daftar kelas berdasarkan nama prodi, dengan pagination (10 per halaman).
// @Description Setiap kelas sudah termasuk deep join: kurikulum → kurikulum_mk (difilter by semester kelas) → mata_kuliah, prodi → jurusan, tahun_akademik, dan daftar mahasiswa.
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter URI tidak valid -> `message: "failed to Get kelas", error: "Key: 'param' Error:..."`
// @Description  - `400` Prodi tidak ditemukan -> `message: "failed to Get kelas", error: "prodi not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to Get kelas", error: "Internal Error"`
// @Tags kelas
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param prodi_name path string true "Nama Prodi" example(teknik-listrik)
// @Param page query int false "Halaman (default 1, 10 per halaman)" example(1)
// @Success      200      {object}  utils.Response[utils.PaginatedData[[]dto.KelasResponse],any]
// @Failure      400      {object}  swagger.ErrGetKelasFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrGetKelasInternalServer
// @Router /api/kelas/prodi/{prodi_name} [get]
func (c *kelasController) GetKelasByProdiName(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.ProdiNameURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var pageQuery utils.PaginationQuery
	if err := ctx.ShouldBindQuery(&pageQuery); err != nil || pageQuery.Page <= 0 {
		pageQuery.Page = 1
	}

	data, total, err := c.kelasService.GetKelasByProdiNamePaginated(ctx.Request.Context(), URI.Name, pageQuery.Page)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	paginated := utils.BuildPaginatedResponse(data, pageQuery.Page, total)
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_KELAS, paginated, path)
	ctx.JSON(http.StatusOK, res)
}

// GetAllKelas godoc
// @Summary      Get All Kelas (Paginated)
// @Description  Mengambil semua daftar kelas dengan pagination (10 per halaman).
// @Description
// @Description  **Akses:** Semua user yang sudah login.
// @Tags         kelas
// @Produce      json
// @Security     ApiKeyAuth
// @Param        page  query  int  false  "Halaman (default 1)"  example(1)
// @Success      200   {object}  utils.Response[utils.PaginatedData[[]dto.KelasResponse],any]
// @Failure      401   {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      500   {object}  swagger.ErrGetKelasInternalServer
// @Router       /api/kelas [get]
func (c *kelasController) GetAllKelas(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var pageQuery utils.PaginationQuery
	if err := ctx.ShouldBindQuery(&pageQuery); err != nil || pageQuery.Page <= 0 {
		pageQuery.Page = 1
	}

	data, total, err := c.kelasService.GetAllKelas(ctx.Request.Context(), pageQuery.Page)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_KELAS, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	paginated := utils.BuildPaginatedResponse(data, pageQuery.Page, total)
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_KELAS, paginated, path)
	ctx.JSON(http.StatusOK, res)
}
