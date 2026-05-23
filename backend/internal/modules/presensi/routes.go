package presensi

import (
	"web-hosting/internal/middlewares"
	"web-hosting/internal/modules/auth/service"
	"web-hosting/internal/modules/presensi/controller"
	"web-hosting/internal/package/constants"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	presensiController := do.MustInvoke[controller.PresensiController](injector)
	jwtService := do.MustInvokeNamed[service.JwtService](injector, constants.JWTService)

	routePresensi := router.Group("/api/presensi")
	routePresensi.Use(middlewares.AuthMiddleware(jwtService))
	{
		routePresensi.POST("/mahasiswa", presensiController.CreatePresensiMahasiswa)
		routePresensi.PUT("/mahasiswa", presensiController.UpdatePresensiMahasiswa)
		routePresensi.PUT("/mahasiswa/qr", middlewares.RoleMiddleware("mahasiswa"), presensiController.UpdatePresensiByQR)
		routePresensi.GET("/mahasiswa", presensiController.GetPresensiMahasiswa)

		routePresensi.POST("/pegawai", presensiController.CreatePresensiPegawai)
		routePresensi.PUT("/pegawai", presensiController.UpdatePresensiPegawai)
		routePresensi.GET("/pegawai", presensiController.GetPresensiPegawai)

	}
}
