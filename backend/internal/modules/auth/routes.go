package auth

import (
	"web-hosting/internal/middlewares"
	"web-hosting/internal/modules/auth/controller"
	"web-hosting/internal/modules/auth/service"
	"web-hosting/internal/package/constants"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	authController := do.MustInvoke[controller.AuthController](injector)
	jwtService := do.MustInvokeNamed[service.JwtService](injector, constants.JWTService)

	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/logout", middlewares.AuthMiddleware(jwtService), authController.Logout)
		authRoutes.POST("/refresh-token", authController.RefreshToken)
		authRoutes.GET("/refresh-token/:refresh_token", authController.FindRefreshToken)
		authRoutes.POST("/reset-password", middlewares.AuthMiddleware(jwtService), authController.ResetPassword)
	}
}
