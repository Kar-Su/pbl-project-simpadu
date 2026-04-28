package auth

import (
	"web-hosting/internal/modules/auth/controller"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	authController := do.MustInvoke[controller.AuthController](injector)
	// Next Invoke JWTAut

	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/logout", authController.Logout)
		authRoutes.POST("/refresh-token", authController.RefreshToken)
		authRoutes.POST("/reset-password", authController.ResetPassword)
	}
}
