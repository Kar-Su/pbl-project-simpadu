package jurusan

import (
	"web-hosting/internal/middlewares"
	"web-hosting/internal/modules/auth/service"
	"web-hosting/internal/modules/jurusan/controller"
	"web-hosting/internal/package/constants"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	jurusanController := do.MustInvoke[controller.JurusanController](injector)
	jwtService := do.MustInvokeNamed[service.JwtService](injector, constants.JWTService)

	jurusanRoute := router.Group("/api/jurusan")
	{
		jurusanRoute.POST("/", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_ADMIN_AKADEMIK, constants.ROLE_SUPER_ADMIN), jurusanController.CreateJurusan)
		jurusanRoute.PUT("/", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_ADMIN_AKADEMIK, constants.ROLE_SUPER_ADMIN), jurusanController.UpdateJurusan)
		jurusanRoute.DELETE("/", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_ADMIN_AKADEMIK, constants.ROLE_SUPER_ADMIN), jurusanController.DeleteJurusan)
		jurusanRoute.GET("/", middlewares.AuthMiddleware(jwtService), jurusanController.GetJurusan)
	}
}
