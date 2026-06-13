package controller

import (
	"net/http"
	"web-hosting/internal/modules/presensi/dto"
	"web-hosting/internal/modules/presensi/service"
	"web-hosting/internal/modules/presensi/validation"
	userService "web-hosting/internal/modules/user/service"
	_ "web-hosting/internal/package/swagger"
	"web-hosting/internal/package/types"
	"web-hosting/internal/package/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type (
	PresensiController interface {
		CreatePresensiMahasiswa(ctx *gin.Context)
		CreatePresensiPegawai(ctx *gin.Context)
		UpdatePresensiMahasiswa(ctx *gin.Context)
		UpdatePresensiPegawai(ctx *gin.Context)
		UpdatePresensiByQR(ctx *gin.Context)
		GetPresensiMahasiswa(ctx *gin.Context)
		GetPresensiPegawai(ctx *gin.Context)
		CountPresensi(ctx *gin.Context)
	}
	presensiController struct {
		db                 *gorm.DB
		presensiValidation *validation.PresensiValidation
		presensiService    service.PresensiService
		userService        userService.UserService
	}
)

func NewPresensiController(injector do.Injector, db *gorm.DB, presensiService service.PresensiService, userService userService.UserService) PresensiController {
	presensiValidation := validation.NewPresensiValidation()
	return &presensiController{
		db:                 db,
		presensiValidation: presensiValidation,
		presensiService:    presensiService,
		userService:        userService,
	}
}

// CreatePresensiMahasiswa godoc
// @Summary Create Presensi Mahasiswa
// @Description Menambahkan presensi baru untuk mahasiswa
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'KelasName' Error:..."`
// @Description  - `400` presensi dengan id tersebut sudah ada -> `message: "failed to create presensi", error: "presensi already exists"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` user tidak memiliki akses -> `message: "kelas anda tidak diizinkan", error: "Forbidden"`
// @Tags presensi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body swagger.PresensiMahasiswaCreateRequest true "Presensi Request"
// @Success      201      {object}  utils.Response[dto.PresensiMahasiswaResponse,any]
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Router /api/presensi/mahasiswa [post]
func (c *presensiController) CreatePresensiMahasiswa(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var req dto.PresensiMahasiswaCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_CREATE_PRESENSI_MAHASISWA, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	data, err := c.presensiService.CreatePresensi(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_CREATE_PRESENSI_MAHASISWA, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.SUCCESS_CREATE_PRESENSI_MAHASISWA, data, path)
	ctx.JSON(http.StatusCreated, res)
}

// CreatePresensiPegawai godoc
// @Summary Create Presensi Pegawai
// @Description Menambahkan presensi baru untuk pegawai
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'PresensiName' Error:..."`
// @Description  - `400` presensi dengan id tersebut sudah ada -> `message: "failed to create presensi", error: "presensi already exists"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` user tidak memiliki akses -> `message: "kelas anda tidak diizinkan", error: "Forbidden"`
// @Tags presensi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success      201      {object}  utils.Response[any,any]
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Router /api/presensi/pegawai [post]
func (c *presensiController) CreatePresensiPegawai(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	if _, err := c.presensiService.CreatePresensi(ctx.Request.Context(), nil); err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_CREATE_PRESENSI_PEGAWAI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.SUCCESS_CREATE_PRESENSI_PEGAWAI, any(nil), path)
	ctx.JSON(http.StatusCreated, res)
}

// UpdatePresensiMahasiswa godoc
// @Summary Update Presensi Mahasiswa
// @Description Update Status presensi untuk mahasiswa
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'KelasName' Error:..."`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` user tidak memiliki akses -> `message: "kelas anda tidak diizinkan", error: "Forbidden"`
// @Tags presensi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body swagger.PresensiMahasiswaUpdateRequest true "Presensi Request"
// @Success      201      {object}  utils.Response[any,any]
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Router /api/presensi/mahasiswa [put]
func (c *presensiController) UpdatePresensiMahasiswa(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	tipePresensi := "mahasiswa"

	var req dto.PresensiMahasiswaUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_UPDATE_DETAIL_PRESENSI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := c.presensiService.UpdatePresensi(ctx.Request.Context(), tipePresensi, req); err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_UPDATE_DETAIL_PRESENSI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.SUCCESS_UPDATE_DETAIL_PRESENSI, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}

// UpdatePresensiPegawai godoc
// @Summary Update Presensi Pegawai
// @Description Update Status presensi untuk pegawai
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'KelasName' Error:..."`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` user tidak memiliki akses -> `message: "kelas anda tidak diizinkan", error: "Forbidden"`
// @Tags presensi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body swagger.PresensiPegawaiUpdateRequest true "Presensi Request"
// @Success      201      {object}  utils.Response[any,any]
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Router /api/presensi/pegawai [put]
func (c *presensiController) UpdatePresensiPegawai(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	tipePresensi := "pegawai"

	var req dto.PresensiPegawaiUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_UPDATE_DETAIL_PRESENSI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := c.presensiService.UpdatePresensi(ctx.Request.Context(), tipePresensi, req); err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_UPDATE_DETAIL_PRESENSI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.SUCCESS_UPDATE_DETAIL_PRESENSI, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}

// UpdatePresensiMahasiswaQR godoc
// @Summary Update Presensi Mahasiswa QR
// @Description Update Status presensi untuk mahasiswa menggunakan QR code
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'KelasName' Error:..."`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` user tidak memiliki akses -> `message: "kelas anda tidak diizinkan", error: "Forbidden"`
// @Tags presensi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param q query dto.UpdatePresensiByQRQuery true "Query Parameters"
// @Success      201      {object}  utils.Response[any,any]
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Router /api/presensi/mahasiswa/qr [put]
func (c *presensiController) UpdatePresensiByQR(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	tipePresensi := "mahasiswa"

	var query dto.UpdatePresensiByQRQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_UPDATE_DETAIL_PRESENSI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userId := uuid.MustParse(ctx.MustGet("user_id").(string))

	user, err := c.userService.GetUserByID(ctx.Request.Context(), userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_UPDATE_DETAIL_PRESENSI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	sesiID, err := uuid.Parse(query.SesiID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_UPDATE_DETAIL_PRESENSI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	req := dto.PresensiMahasiswaUpdateRequest{
		PresensiID: sesiID,
		Detail: []dto.DetailPresensiUpdateRequest{
			{DetailID: *user.DetailId, Status: "hadir"},
		},
	}

	if err := c.presensiService.UpdatePresensi(ctx.Request.Context(), tipePresensi, req); err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_UPDATE_DETAIL_PRESENSI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.SUCCESS_UPDATE_DETAIL_PRESENSI, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}

// GetPresensiMahasiswa godoc
// @Summary Get Presensi Mahasiswa
// @Description Mengambil presensi mahasiswa berdasarkan id
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Tidak ada presensi dengan id tersebut -> `message: "failed to get presensi", error: "presensi not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` user tidak memiliki akses -> `message: "kelas anda tidak diizinkan", error: "Forbidden"`
// @Tags presensi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sesi_id query dto.GetPresensiMahasiswaQuery true "ID sesi"
// @Success      200      {object}  utils.Response[dto.PresensiMahasiswaResponse,any]
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Router /api/presensi/mahasiswa [get]
func (c *presensiController) GetPresensiMahasiswa(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	tipePresensi := "mahasiswa"

	var query dto.GetPresensiMahasiswaQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_GET_PRESENSI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	data, err := c.presensiService.GetPresensi(ctx.Request.Context(), tipePresensi, query.ID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_GET_PRESENSI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.SUCCESS_GET_PRESENSI, data, path)
	ctx.JSON(http.StatusOK, res)
}

// GetPresensiPegawai godoc
// @Summary Get Presensi Pegawai
// @Description Mengambil presensi pegawai berdasarkan id
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Tidak ada presensi dengan id tersebut -> `message: "failed to get presensi", error: "presensi not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` user tidak memiliki akses -> `message: "kelas anda tidak diizinkan", error: "Forbidden"`
// @Tags presensi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param q query dto.GetPresensiPegawaiQuery true "presensi ID or date (YYYY-MM-DD)"
// @Param page query int false "Halaman (default 1, 10 per halaman)" example(1)
// @Success      200      {object}  utils.Response[any,any]
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Router /api/presensi/pegawai [get]
func (c *presensiController) GetPresensiPegawai(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	tipePresensi := "pegawai"

	var query dto.GetPresensiPegawaiQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_GET_PRESENSI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var filter any

	if query.Date != "" {
		var err error
		filter, err = types.ParseString(query.Date)
		if err != nil {
			res := utils.BuildResponseFailed(dto.FAILED_GET_PRESENSI, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}

	} else if query.ID != "" {
		var err error
		filter, err = uuid.Parse(query.ID)
		if err != nil {
			res := utils.BuildResponseFailed(dto.FAILED_GET_PRESENSI, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}
	} else {
		var pageQuery utils.PaginationQuery
		if err := ctx.ShouldBindQuery(&pageQuery); err != nil || pageQuery.Page <= 0 {
			pageQuery.Page = 1
		}

		data, total, err := c.presensiService.GetAllPresensiPaginated(ctx.Request.Context(), tipePresensi, nil, pageQuery.Page)
		if err != nil {
			res := utils.BuildResponseFailed(dto.FAILED_GET_PRESENSI, err.Error(), nil, path)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}

		paginated := utils.BuildPaginatedResponse(data, pageQuery.Page, total)
		res := utils.BuildResponseSuccess(dto.SUCCESS_GET_PRESENSI, paginated, path)
		ctx.JSON(http.StatusOK, res)
		return
	}

	data, err := c.presensiService.GetPresensi(ctx.Request.Context(), tipePresensi, filter)
	if err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_GET_PRESENSI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.SUCCESS_GET_PRESENSI, data, path)
	ctx.JSON(http.StatusOK, res)
}

// CountPresensi godoc
// @Summary Count Presensi
// @Description Menghitung jumlah presensi pegawai berdasarkan tipe presensi
// @Description
// @Description  **Akses:** Logged User
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Tidak ada presensi dengan id tersebut -> `message: "failed to count presensi", error: "presensi not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` user tidak memiliki akses -> `message: "kelas anda tidak diizinkan", error: "Forbidden"`
// @Tags presensi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param tipe query dto.TipePresensiQuery true "Tipe presensi (mahasiswa/pegawai)" example(mahasiswa)
// @Success      200      {object}  utils.Response[any,any]
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Router /api/presensi/count [get]
func (c *presensiController) CountPresensi(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var query dto.TipePresensiQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_COUNT_PRESENSI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	count, err := c.presensiService.CountPresensi(ctx.Request.Context(), query.Tipe)
	if err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_COUNT_PRESENSI, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.SUCCESS_COUNT_PRESENSI, count, path)
	ctx.JSON(http.StatusOK, res)
}
