package khs

import (
	"web-hosting/internal/middlewares"
	"web-hosting/internal/modules/auth/service"
	"web-hosting/internal/modules/khs/controller"
	"web-hosting/internal/package/constants"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	khsController := do.MustInvoke[controller.KHSController](injector)
	jwtService := do.MustInvokeNamed[service.JwtService](injector, constants.JWTService)

	khsRoute := router.Group("/api/khs")
	{
		khsRoute.POST("/nilai", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_PEGAWAI, constants.ROLE_DOSEN), khsController.Create)
		khsRoute.GET("/", middlewares.AuthMiddleware(jwtService), khsController.GetKHS)
	}
}
