package user

import (
	"web-hosting/internal/modules/user/controller"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	userController := do.MustInvoke[controller.UserController](injector)
	// Next Harus Invoce JWTAut
	userRoutes := router.Group("/api/user")
	{
		userRoutes.GET("/:id", userController.GetUser)
		userRoutes.PUT("/super/:id", userController.UpdateAdmin)
		userRoutes.DELETE("/super/:id", userController.DeleteAdmin)

		userRoutes.GET("/sync/:role_name/:detail_id", userController.GetUserNonAdmin)
		userRoutes.PUT("/sync/:role_name/:detail_id", userController.UpdateNonAdmin)
		userRoutes.DELETE("/sync/:role_name/:detail_id", userController.DeleteNonAdmin)
	}
}
