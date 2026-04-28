package role

import (
	"web-hosting/internal/modules/role/controller"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	roleController := do.MustInvoke[controller.RoleController](injector)
	//Next AuthMiddlewarenya
	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/role", roleController.GetAllRole)

		apiRoutes.PUT("/super/role", roleController.Update)
		apiRoutes.POST("/super/role", roleController.Create)
		apiRoutes.DELETE("/super/role/:id", roleController.Delete)

	}

}
