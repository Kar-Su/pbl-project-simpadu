package pengampu

import (
	"web-hosting/internal/middlewares"
	"web-hosting/internal/modules/auth/service"
	"web-hosting/internal/modules/pengampu/controller"
	"web-hosting/internal/package/constants"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(r *gin.Engine, injector do.Injector) {
	pengampuController := do.MustInvoke[controller.PengampuController](injector)
	jwtService := do.MustInvokeNamed[service.JwtService](injector, constants.JWTService)

	pengampuRoute := r.Group("/api/pengampu")
	pengampuRoute.Use(middlewares.AuthMiddleware(jwtService))
	{
		pengampuRoute.POST("", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_PEGAWAI), pengampuController.CreatePengampu)
		pengampuRoute.PUT("/:pengampu_id", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_PEGAWAI), pengampuController.UpdatePengampuByID)
		pengampuRoute.DELETE("/:pengampu_id", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_PEGAWAI), pengampuController.DeletePengampuByID)
		pengampuRoute.GET("/:pengampu_id", pengampuController.GetPengampuByID)
		pengampuRoute.GET("/kelas/:kelas_id", pengampuController.GetPengampuByKelasID)
	}
}
