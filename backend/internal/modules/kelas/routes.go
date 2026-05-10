package kelas

import (
	"web-hosting/internal/middlewares"
	"web-hosting/internal/modules/auth/service"
	"web-hosting/internal/modules/kelas/controller"
	"web-hosting/internal/package/constants"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	kelasController := do.MustInvoke[controller.KelasController](injector)
	pivotController := do.MustInvoke[controller.KelasMahasiswaController](injector)
	jwtService := do.MustInvokeNamed[service.JwtService](injector, constants.JWTService)

	kelasRoute := router.Group("/api/kelas")
	kelasRoute.Use(middlewares.AuthMiddleware(jwtService))
	{
		kelasRoute.POST("/", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), kelasController.CreateKelas)
		kelasRoute.PUT("/:kelas_id", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), kelasController.UpdateKelas)
		kelasRoute.DELETE("/:kelas_id", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), kelasController.DeleteKelas)
		kelasRoute.GET("/:kelas_id", kelasController.GetKelasByID)
		kelasRoute.GET("/prodi/:prodi_name", kelasController.GetKelasByProdiName)

		kelasRoute.POST("/mahasiswa", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), pivotController.AssignMahasiswaToKelas)
		kelasRoute.DELETE("/:kelas_id/mahasiswa/:mahasiswa_id", middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK), pivotController.RemoveMahasiswaFromKelas)
		kelasRoute.GET("/:kelas_id/mahasiswa", pivotController.GetMahasiswaByKelas)
		kelasRoute.GET("/mahasiswa/:mahasiswa_id", pivotController.GetAllKelasMahasiswa)
	}
}
