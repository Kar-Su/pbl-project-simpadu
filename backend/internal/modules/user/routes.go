package user

import (
	"web-hosting/internal/middlewares"
	"web-hosting/internal/modules/auth/service"
	"web-hosting/internal/modules/user/controller"
	"web-hosting/internal/package/constants"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	userController := do.MustInvoke[controller.UserController](injector)
	jwtService := do.MustInvokeNamed[service.JwtService](injector, constants.JWTService)
	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/me", middlewares.AuthMiddleware(jwtService), userController.Me)
		apiRoutes.GET("/user/:id", middlewares.AuthMiddleware(jwtService), userController.GetUser)
		apiRoutes.POST("/super/user", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN), userController.RegisterAdmin)
		apiRoutes.PUT("/super/user/:id", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN), userController.UpdateAdmin)
		apiRoutes.DELETE("/super/user/:id", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN), userController.DeleteAdmin)

		apiRoutes.GET("/user/role/:role_name", middlewares.AuthMiddleware(jwtService), userController.GetUserByRole)
		apiRoutes.GET("/user/email/:email", middlewares.AuthMiddleware(jwtService), userController.GetUserByEmail)
		apiRoutes.POST("/user", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN), userController.RegisterNonAdmin)

		apiRoutes.GET("/user/sync/:role_name/:detail_id", middlewares.AuthMiddleware(jwtService), userController.GetUserNonAdmin)
		apiRoutes.PUT("/user/sync/:role_name/:detail_id", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN), userController.UpdateNonAdmin)
		apiRoutes.DELETE("/user/sync/:role_name/:detail_id", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN), userController.DeleteNonAdmin)
	}
}
