package controller

import (
	"net/http"
	"web-hosting/internal/modules/khs/dto"
	"web-hosting/internal/modules/khs/service"
	_ "web-hosting/internal/package/swagger"
	"web-hosting/internal/package/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type KHSController interface {
	Create(ctx *gin.Context)
	GetKHS(ctx *gin.Context)
}

type khsController struct {
	db         *gorm.DB
	khsService service.KhsService
}

func NewKHSController(injector do.Injector, db *gorm.DB, khsService service.KhsService) KHSController {
	return &khsController{
		db:         db,
		khsService: khsService,
	}
}

// CreateNilai godoc
// @Summary Create Nilai Baru
// @Description Menambahkan nilai baru ke sistem
// @Description
// @Description  **Akses:** Admin Akademik, Dosen
// @Tags khs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.CreateKhsRequest true "nilai Request"
// @Success      201      {object}  utils.Response[any,any]
// @Failure      400      {object}  swagger.ErrCreateKurikulumFailed
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Router /api/khs/nilai [post]
func (c *khsController) Create(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var req dto.CreateKhsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed("Failed Create KHS/Nilai", err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.khsService.Create(ctx, &req)
	if err != nil {
		res := utils.BuildResponseFailed("Failed Create KHS/Nilai", err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess("Success Create KHS/Nilai", any(nil), path)
	ctx.JSON(http.StatusCreated, res)
}

// GetKhs godoc
// @Summary      Get khs (Paginated)
// @Description  Mengambil daftar KHS dengan pagination (10 per halaman).
// @Description  Logged in user.
// @Description
// @Description  **Akses:** Semua user yang sudah login (Authenticated User).
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter URI tidak valid -> `message: "bad request", error: "Key: 'KhsId' Error:..."`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Tags         khs
// @Produce      json
// @Security     ApiKeyAuth
// @Param        semester  query  int  false  "Semester"  example(1)
// @Param        prodi_name  query  string  false  "Nama Prodi"  example("Teknik Informatika")
// @Param        page       query  int  false  "Halaman (default 1, 10 per halaman)"  example(1)
// @Success      200        {object}  utils.Response[utils.PaginatedData[[]dto.KHSResponse],any]
// @Failure      401        {object}  swagger.ErrUnauthorizedInvalidToken
// @Router       /api/khs [get]
func (c *khsController) GetKHS(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var filter dto.FilterQuery
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		res := utils.BuildResponseFailed("Failed Get KHS/Nilai", err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var paginated utils.PaginationQuery
	if err := ctx.ShouldBindQuery(&paginated); err != nil {
		res := utils.BuildResponseFailed("Failed Get KHS/Nilai", err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	khs, err := c.khsService.GetKHS(ctx, &filter, paginated.Page)
	if err != nil {
		res := utils.BuildResponseFailed("Failed Get KHS/Nilai", err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess("Success Get KHS/Nilai", khs, path)
	ctx.JSON(http.StatusOK, res)
}
