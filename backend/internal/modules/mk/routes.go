package mk

import (
	"web-hosting/internal/middlewares"
	"web-hosting/internal/modules/auth/service"
	"web-hosting/internal/modules/mk/controller"
	"web-hosting/internal/package/constants"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	mkController := do.MustInvoke[controller.MkController](injector)
	jwtService := do.MustInvokeNamed[service.JwtService](injector, constants.JWTService)

	mkRoutes := router.Group("/api/mata-kuliah")
	{
		mkRoutes.POST("", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), mkController.CreateMk)
		mkRoutes.PUT("/:kode", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), mkController.UpdateMk)
		mkRoutes.DELETE("/:kode", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), mkController.DeleteMk)
		mkRoutes.GET("", middlewares.AuthMiddleware(jwtService), mkController.GetMk)
	}
}
