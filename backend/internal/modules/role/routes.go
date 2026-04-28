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
	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/role", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN), roleController.GetAllRole)

		apiRoutes.PUT("/super/role", roleController.Update, middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN))
		apiRoutes.POST("/super/role", roleController.Create, middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN))
		apiRoutes.DELETE("/super/role/:id", roleController.Delete, middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN))

	}

}
