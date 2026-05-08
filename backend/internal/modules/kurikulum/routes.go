package kurikulum

import (
	"web-hosting/internal/middlewares"
	"web-hosting/internal/modules/auth/service"
	"web-hosting/internal/modules/kurikulum/controller"
	"web-hosting/internal/package/constants"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(route *gin.Engine, injector do.Injector) {
	kurikulumController := do.MustInvoke[controller.KurikulumController](injector)
	pivotController := do.MustInvoke[controller.PivotController](injector)
	jwtService := do.MustInvokeNamed[service.JwtService](injector, constants.JWTService)

	kurikulumRoute := route.Group("/api/kurikulum")
	kurikulumRoute.Use(middlewares.AuthMiddleware(jwtService))
	{
		kurikulumRoute.POST("/", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), kurikulumController.CreateKurikulum)
		kurikulumRoute.PUT("/", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), kurikulumController.UpdateKurikulum)
		kurikulumRoute.DELETE("/", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), kurikulumController.DeleteKurikulum)
		kurikulumRoute.GET("/", kurikulumController.GetKurikulum)

		kurikulumRoute.POST("/mata-kuliah", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), pivotController.CreatePivot)
		kurikulumRoute.PUT("/:kurikulum_kode/mata-kuliah/:mk_kode", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), pivotController.UpdatePivot)
		kurikulumRoute.DELETE("/:kurikulum_kode/mata-kuliah/:mk_kode", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), pivotController.DeletePivot)
	}
}
