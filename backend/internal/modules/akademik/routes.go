package akademik

import (
	"web-hosting/internal/middlewares"
	"web-hosting/internal/modules/akademik/controller"
	"web-hosting/internal/modules/auth/service"
	"web-hosting/internal/package/constants"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(routes *gin.Engine, injector do.Injector) {
	akademikController := do.MustInvoke[controller.TahunAkademikController](injector)
	jwtService := do.MustInvokeNamed[service.JwtService](injector, constants.JWTService)

	routeAkademik := routes.Group("/api/tahun-akademik")
	{
		routeAkademik.Use(middlewares.AuthMiddleware(jwtService))
		routeAkademik.POST("/", middlewares.RoleMiddleware(constants.ROLE_ADMIN_AKADEMIK, constants.ROLE_SUPER_ADMIN), akademikController.CreateTahunAkademik)
		routeAkademik.PUT("/:id", middlewares.RoleMiddleware(constants.ROLE_ADMIN_AKADEMIK, constants.ROLE_SUPER_ADMIN), akademikController.UpdateTahunAkademik)
		routeAkademik.DELETE("/:id", middlewares.RoleMiddleware(constants.ROLE_ADMIN_AKADEMIK, constants.ROLE_SUPER_ADMIN), akademikController.DeleteTahunAkademik)
		routeAkademik.GET("/:id", akademikController.GetTahunAkademik)
		routeAkademik.GET("/", akademikController.GetAllTahunAkademik)
		routeAkademik.GET("/status/:status", akademikController.GetTahunAkademikByStatus)
	}
}
