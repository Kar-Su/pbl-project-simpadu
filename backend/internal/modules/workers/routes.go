package workers

import (
	"web-hosting/internal/middlewares"
	"web-hosting/internal/modules/auth/service"
	"web-hosting/internal/modules/workers/controller"
	"web-hosting/internal/package/constants"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	workerController := do.MustInvoke[controller.WorkerController](injector)
	jwtService := do.MustInvokeNamed[service.JwtService](injector, constants.JWTService)

	routePresensi := router.Group("/api/workers")
	routePresensi.GET("/presensi/status", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_PEGAWAI), workerController.GetPresensiStatus)
	routePresensi.POST("/presensi/start", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_PEGAWAI), workerController.StartPresensi)
	routePresensi.POST("/presensi/stop", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_PEGAWAI), workerController.StopPresensi)
}
