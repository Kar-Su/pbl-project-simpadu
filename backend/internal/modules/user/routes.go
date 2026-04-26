package user

import (
	"web-hosting/internal/modules/user/controller"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	userController := do.MustInvoke[controller.UserController](injector)
	// Next Harus Invoce JWTAut
	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/user/:id", userController.GetUser)
		apiRoutes.POST("/super/user", userController.RegisterAdmin)
		apiRoutes.PUT("/super/user/:id", userController.UpdateAdmin)
		apiRoutes.DELETE("/super/user/:id", userController.DeleteAdmin)

		apiRoutes.GET("/user/role/:role_name", userController.GetUserByRole)
		apiRoutes.GET("/user/email/:email", userController.GetUserByEmail)
		apiRoutes.POST("/user", userController.RegisterNonAdmin)

		apiRoutes.GET("/user/sync/:role_name/:detail_id", userController.GetUserNonAdmin)
		apiRoutes.PUT("/user/sync/:role_name/:detail_id", userController.UpdateNonAdmin)
		apiRoutes.DELETE("/user/sync/:role_name/:detail_id", userController.DeleteNonAdmin)
	}
}
