package role

import (
	"web-hosting/internal/middlewares"
	"web-hosting/internal/modules/auth/service"
	"web-hosting/internal/modules/role/controller"
	"web-hosting/internal/package/constants"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	roleController := do.MustInvoke[controller.RoleController](injector)
	jwtService := do.MustInvokeNamed[service.JwtService](injector, constants.JWTService)

	rolesRoute := router.Group("/api/roles")
	{
		rolesRoute.GET("", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK, constants.ROLE_ADMIN_KEUANGAN, constants.ROLE_ADMIN_MAHASISWA), roleController.GetAllRole)
		rolesRoute.POST("", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN), roleController.Create)
		rolesRoute.PUT("/:role_name", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN), roleController.Update)
		rolesRoute.DELETE("/:role_name", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN), roleController.Delete)
	}
}
