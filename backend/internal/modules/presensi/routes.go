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
		routePresensi.POST("/mahasiswa", middlewares.AuthMiddleware(jwtService), presensiController.CreatePresensiMahasiswa)
		routePresensi.PUT("/mahasiswa", middlewares.AuthMiddleware(jwtService), presensiController.UpdatePresensiMahasiswa)
		routePresensi.PUT("/mahasiswa/qr", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware("mahasiswa"), presensiController.UpdatePresensiByQR)
		routePresensi.GET("/mahasiswa", middlewares.AuthMiddleware(jwtService), presensiController.GetPresensiMahasiswa)

		routePresensi.POST("/pegawai", middlewares.AuthMiddleware(jwtService), presensiController.CreatePresensiPegawai)
		routePresensi.PUT("/pegawai", middlewares.AuthMiddleware(jwtService), presensiController.UpdatePresensiPegawai)
		routePresensi.GET("/pegawai", middlewares.AuthMiddleware(jwtService), presensiController.GetPresensiPegawai)

		routePresensi.GET("/count", middlewares.AuthMiddleware(jwtService), presensiController.CountPresensi)

	}
}
