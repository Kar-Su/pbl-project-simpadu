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
		// CRUD kurikulum tanpa param
		kurikulumRoute.POST("", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), kurikulumController.CreateKurikulum)
		kurikulumRoute.GET("", kurikulumController.GetKurikulum)

		// Pivot CREATE: kurikulum_kode dari body request
		kurikulumRoute.POST("/mata-kuliah", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), pivotController.CreatePivot)

		// CRUD kurikulum by kode + sub-resource pivot
		// Catatan: Gin mengharuskan semua wildcard pada level yang sama memiliki nama yang identik.
		// /:kode digunakan konsisten untuk semua operasi agar tidak ada konflik.
		kurikulumRoute.GET("/:kode", kurikulumController.GetKurikulum)
		kurikulumRoute.PUT("/:kode", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), kurikulumController.UpdateKurikulum)
		kurikulumRoute.DELETE("/:kode", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), kurikulumController.DeleteKurikulum)

		// Sub-resource pivot: wildcard harus bernama sama (:kode) seperti parent route
		kurikulumRoute.PUT("/:kode/mata-kuliah/:mk_kode", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), pivotController.UpdatePivot)
		kurikulumRoute.DELETE("/:kode/mata-kuliah/:mk_kode", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), pivotController.DeletePivot)
	}
}
