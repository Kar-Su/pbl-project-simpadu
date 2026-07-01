package controller

import (
	"net/http"
	"web-hosting/internal/workers"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type WorkerController interface {
	StartPresensi(ctx *gin.Context)
	StopPresensi(ctx *gin.Context)
	GetPresensiStatus(ctx *gin.Context)
}

type workerController struct {
	scheduler workers.Schedule
}

func NewWorkerController(injector do.Injector) WorkerController {
	scheduler := do.MustInvoke[workers.Schedule](injector)
	return &workerController{scheduler: scheduler}
}

// StartWorkersPresensi godoc
// @Summary      Otomatis Start Presensi pegawai
// @Description
// @Description  **Akses:** Admin Pegawai.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `400` Worker presensi already on"`
// @Tags         workers
// @Produce      json
// @Security     ApiKeyAuth
// @Router       /api/workers/presensi/start [post]
func (c *workerController) StartPresensi(ctx *gin.Context) {
	if c.scheduler.IsPresensiRunning() {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Worker presensi already on"})
		return
	}

	c.scheduler.StartPresensi()
	ctx.JSON(http.StatusOK, gin.H{"message": "Worker presensi on"})
}

// StopWorkersPresensi godoc
// @Summary      Otomatis Stop Presensi pegawai
// @Description
// @Description  **Akses:** Admin Pegawai.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Description  - `400` message": "Worker presensi already off"`
// @Tags         workers
// @Produce      json
// @Security     ApiKeyAuth
// @Router       /api/workers/presensi/stop [post]
func (c *workerController) StopPresensi(ctx *gin.Context) {
	if !c.scheduler.IsPresensiRunning() {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Worker presensi already off"})
		return
	}

	c.scheduler.StopPresensi()
	ctx.JSON(http.StatusOK, gin.H{"message": "Worker presensi off"})
}

// GetStatusWorkerPresensi godoc
// @Summary      Get Status Presensi pegawai
// @Description
// @Description  **Akses:** Admin Pegawai.
// @Description
// @Description  **Error yang mungkin terjadi:**
// @Description  - `401` Authorization header tidak ada -> `message: "failed_auth", error: "Authorization header missing"`
// @Description  - `401` Format header salah (bukan "Bearer ...") -> `message: "failed_auth", error: "invalid authentication header"`
// @Description  - `401` Token JWT tidak valid atau kedaluwarsa -> `message: "failed_auth", error: "invalid token"`
// @Tags         workers
// @Produce      json
// @Security     ApiKeyAuth
// @Router       /api/workers/presensi/status [get]
func (c *workerController) GetPresensiStatus(ctx *gin.Context) {
	status := "off"
	if c.scheduler.IsPresensiRunning() {
		status = "on"
	}
	ctx.JSON(http.StatusOK, gin.H{"status": status})
}
