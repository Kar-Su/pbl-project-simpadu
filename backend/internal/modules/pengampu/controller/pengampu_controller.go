package controller

import (
	"errors"
	"net/http"
	"web-hosting/internal/modules/pengampu/dto"
	"web-hosting/internal/modules/pengampu/service"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/utils"

	_ "web-hosting/internal/package/swagger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type (
	PengampuController interface {
		CreatePengampu(ctx *gin.Context)
		UpdatePengampuByID(ctx *gin.Context)
		DeletePengampuByID(ctx *gin.Context)
		GetPengampuByID(ctx *gin.Context)
		GetPengampuByKelasID(ctx *gin.Context)
		GetPengampuByDosenID(ctx *gin.Context)
	}

	pengampuController struct {
		db              *gorm.DB
		pengampuService service.PengampuService
	}
)

func NewPengampuController(injector do.Injector, db *gorm.DB, pengampuService service.PengampuService) PengampuController {
	return &pengampuController{
		db:              db,
		pengampuService: pengampuService,
	}
}

// CreatePengampu godoc
// @Summary      Buat Pengampu Baru
// @Description  Menambahkan pengampu baru ke dalam sistem.
// @Description
// @Description  **Akses:** Super Admin, Admin Pegawai.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'Pengampu...' Error:..."`
// @Description  - `400` pengampu dengan nama tersebut sudah ada -> `message: "failed to create pengampu", error: "pengampu already exists"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` Role user tidak memiliki akses -> `message: "Role anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to create role", error: "Internal Error"`
// @Tags         pengampu
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      swagger.CreatePengampuRequest  true  "Payload Create Pengampu"
// @Success      201      {object}  utils.Response[any, any]
// @Failure      401      {object}  swagger.ErrUnauthorizedInvalidToken
// @Failure      403      {object}  swagger.ErrForbiddenAccess
// @Failure      500      {object}  swagger.ErrCreateRoleInternalServer
// @Router       /api/pengampu [post]
func (c *pengampuController) CreatePengampu(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var req dto.CreatePengampuRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_CREATE_PENGAMPU, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := c.pengampuService.CreatePengampu(ctx.Request.Context(), req); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.FAILED_CREATE_PENGAMPU, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.SUCCESS_CREATE_PENGAMPU, any(nil), path)
	ctx.JSON(http.StatusCreated, res)
}

// UpdatePengampuByID godoc
// @Summary      Update Pengampu
// @Description  Mengubah data pengampu berdasarkan ID.
// @Description
// @Description  **Akses:** Super Admin, Admin Pegawai.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter URI tidak valid -> `message: "failed to validate pengampu uri", error: "Key: 'pengampuID..' Error:..."`
// @Description  - `400` Body tidak valid / field wajib kosong -> `message: "failed to get request", error: "Key: 'Pengampu..' Error:..."`
// @Description  - `400` Pengampu dengan ID tersebut tidak ditemukan -> `message: "failed to update pengampu", error: "pengampu not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` Role user tidak memiliki akses -> `message: "Role anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to update role", error: "Internal Error"`
// @Tags         pengampu
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        pengampu_id  path      string                 true  "ID Pengampu yang Akan Diubah"  example(rrh55dfb-h4o796gntyt)
// @Param        request    body      swagger.UpdatePengampuRequest  true  "Payload Update Pengampu"
// @Success      200        {object}  utils.Response[dto.PengampuResponse,any]
// @Failure      500      {object}  swagger.ErrCreateRoleInternalServer
// @Router       /api/pengampu/{pengampu_id} [put]
func (c *pengampuController) UpdatePengampuByID(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.PengampuIdURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_UPDATE_PENGAMPU, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var req dto.UpdatePengampuRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_UPDATE_PENGAMPU, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	data, err := c.pengampuService.UpdatePengampuByID(ctx.Request.Context(), uuid.MustParse(URI.PengampuID), req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.FAILED_UPDATE_PENGAMPU, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.SUCCESS_UPDATE_PENGAMPU, data, path)
	ctx.JSON(http.StatusOK, res)
}

// DeletePengampuById godoc
// @Summary      Hapus Pengampu by id
// @Description  Menghapus pengampu dari sistem secara permanen berdasarkan id.
// @Description
// @Description  **Akses:** Khusus Super Admin.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter URI tidak valid -> `message: "failed to validate role uri", error: "Key: 'RoleName' Error:..."`
// @Description  - `400` Pengampu dengan id tersebut tidak ditemukan -> `message: "failed to delete pengampu", error: "pengampu not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `403` Role user tidak memiliki akses -> `message: "Role anda tidak diizinkan", error: "Forbidden"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to delete pengampu", error: "Internal Error"`
// @Tags         pengampu
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        pengampu_id  path      string  true  "ID Pengampu yang Akan Dihapus"  example(g202gdvs-h4rtrh-56tjt)
// @Success      200        {object}  utils.Response[any,any]
// @Failure      500      {object}  swagger.ErrCreateRoleInternalServer
// @Router       /api/pengampu/{pengampu_id} [delete]
func (c *pengampuController) DeletePengampuByID(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.PengampuIdURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_DELETE_PENGAMPU, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := c.pengampuService.DeletePengampuByID(ctx.Request.Context(), uuid.MustParse(URI.PengampuID)); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.FAILED_DELETE_PENGAMPU, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.SUCCESS_DELETE_PENGAMPU, any(nil), path)
	ctx.JSON(http.StatusOK, res)
}

// GetPengampuById godoc
// @Summary      Get Pengampu by id
// @Description  Mendapatkan pengampu dari sistem berdasarkan id.
// @Description
// @Description  **Akses:** Logged User.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter URI tidak valid -> `message: "failed to validate role uri", error: "Key: 'RoleName' Error:..."`
// @Description  - `400` Pengampu dengan id tersebut tidak ditemukan -> `message: "failed to delete pengampu", error: "pengampu not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to delete pengampu", error: "Internal Error"`
// @Tags         pengampu
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        pengampu_id  path      string  true  "ID Pengampu"  example(32kg-k3943-)
// @Success      200        {object}  utils.Response[dto.PengampuResponse,any]
// @Failure      500        {object}  swagger.ErrDeleteRoleInternalServer
// @Router       /api/pengampu/{pengampu_id} [get]
func (c *pengampuController) GetPengampuByID(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.PengampuIdURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_GET_PENGAMPU, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	pengampuID, err := uuid.Parse(URI.PengampuID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_GET_PENGAMPU, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	data, err := c.pengampuService.GetPengampuByID(ctx.Request.Context(), pengampuID)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.FAILED_GET_PENGAMPU, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.SUCCESS_GET_PENGAMPU, data, path)
	ctx.JSON(http.StatusOK, res)
}

// GetPengampuByKelas godoc
// @Summary      Get Pengampu by kelas id
// @Description  Mendapatkan pengampu dari sistem berdasarkan id.
// @Description
// @Description  **Akses:** Logged User.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter URI tidak valid -> `message: "failed to validate role uri", error: "Key: 'RoleName' Error:..."`
// @Description  - `400` Pengampu dengan id tersebut tidak ditemukan -> `message: "failed to delete pengampu", error: "pengampu not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to delete pengampu", error: "Internal Error"`
// @Tags         pengampu
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        kelas_id  path      string  true  "UUID kelas"  example(da322-f33)
// @Success      200        {object}  utils.Response[[]dto.PengampuResponse,any]
// @Failure      500        {object}  swagger.ErrDeleteRoleInternalServer
// @Router       /api/pengampu/kelas/{kelas_id} [get]
func (c *pengampuController) GetPengampuByKelasID(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.PengampuKelasIdURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_GET_PENGAMPU, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	kelasID, err := uuid.Parse(URI.KelasID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_GET_PENGAMPU, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	data, err := c.pengampuService.GetPengampuByKelasID(ctx.Request.Context(), kelasID)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.FAILED_GET_PENGAMPU, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.SUCCESS_GET_PENGAMPU, data, path)
	ctx.JSON(http.StatusOK, res)
}

// GetPengampuByDosenID godoc
// @Summary      Get Pengampu by dosen id
// @Description  Mendapatkan pengampu dari sistem berdasarkan id.
// @Description
// @Description  **Akses:** Logged User.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `400` Parameter URI tidak valid -> `message: "failed to validate role uri", error: "Key: 'PengampuID' Error:..."`
// @Description  - `400` Pengampu dengan id tersebut tidak ditemukan -> `message: "failed to delete pengampu", error: "pengampu not found"`
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `500` Kesalahan internal server -> `message: "failed to delete pengampu", error: "Internal Error"`
// @Tags         pengampu
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        dosen_id  path      string  true  "UUID dosen"  example(da322-f33)
// @Success      200        {object}  utils.Response[[]dto.PengampuResponse,any]
// @Failure      500        {object}  swagger.ErrDeleteRoleInternalServer
// @Router       /api/pengampu/dosen/{dosen_id} [get]
func (c *pengampuController) GetPengampuByDosenID(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var URI dto.PengampuDosenIdURI
	if err := ctx.ShouldBindUri(&URI); err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_GET_PENGAMPU, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	dosenID, err := uuid.Parse(URI.DosenID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.FAILED_GET_PENGAMPU, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	data, err := c.pengampuService.GetPengampuByDosenID(ctx.Request.Context(), dosenID)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, constants.ErrInternalErr) {
			status = http.StatusInternalServerError
		}
		res := utils.BuildResponseFailed(dto.FAILED_GET_PENGAMPU, err.Error(), nil, path)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.SUCCESS_GET_PENGAMPU, data, path)
	ctx.JSON(http.StatusOK, res)
}
