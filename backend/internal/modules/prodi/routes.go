package prodi

import (
	"web-hosting/internal/middlewares"
	"web-hosting/internal/modules/auth/service"
	"web-hosting/internal/modules/prodi/controller"
	"web-hosting/internal/package/constants"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	prodiController := do.MustInvoke[controller.ProdiController](injector)
	jwtService := do.MustInvokeNamed[service.JwtService](injector, constants.JWTService)

	prodiRoute := router.Group("/api/prodi")
	{
		prodiRoute.POST("/", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), prodiController.CreateProdi)
		prodiRoute.PUT("/", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), prodiController.UpdateProdi)
		prodiRoute.DELETE("/", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), prodiController.DeleteProdi)
		prodiRoute.GET("/", middlewares.AuthMiddleware(jwtService), prodiController.GetProdi)
		prodiRoute.GET("/jurusan/:jurusan_name", middlewares.AuthMiddleware(jwtService), prodiController.GetProdiByJurusan)
	}
}
