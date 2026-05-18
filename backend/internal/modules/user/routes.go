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
		// Current user
		apiRoutes.GET("/me", middlewares.AuthMiddleware(jwtService), userController.Me)

		// Admin user CRUD - resource: /api/users
		apiRoutes.GET("/users", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_AKADEMIK, constants.ROLE_ADMIN_PEGAWAI, constants.ROLE_ADMIN_MAHASISWA, constants.ROLE_ADMIN_KEUANGAN), userController.GetAllUsers)
		apiRoutes.GET("/users/:id", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN), userController.GetUser)
		apiRoutes.POST("/users/admins", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN), userController.RegisterAdmin)
		apiRoutes.PUT("/users/:id", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN), userController.UpdateAdmin)
		apiRoutes.DELETE("/users/:id", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN), userController.DeleteAdmin)

		// Non-admin user CRUD
		apiRoutes.GET("/users/search", middlewares.AuthMiddleware(jwtService), userController.GetUserByEmail)
		apiRoutes.POST("/users", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_MAHASISWA, constants.ROLE_ADMIN_PEGAWAI), userController.RegisterNonAdmin)
		apiRoutes.GET("/users/count", middlewares.AuthMiddleware(jwtService), userController.CountAllUsers)

		// Filter by role: GET /api/users/roles/:role_name?page=N
		apiRoutes.GET("/users/roles/:role_name", middlewares.AuthMiddleware(jwtService), userController.GetUserByRole)

		// Sync endpoints: untuk sinkronisasi data user dengan service lain
		apiRoutes.GET("/users/sync/:role_name/:detail_id", middlewares.AuthMiddleware(jwtService), userController.GetUserNonAdmin)
		apiRoutes.PUT("/users/sync/:role_name/:detail_id", middlewares.AuthMiddleware(jwtService), userController.UpdateNonAdmin)
		apiRoutes.DELETE("/users/sync/:role_name/:detail_id", middlewares.AuthMiddleware(jwtService), middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_PEGAWAI, constants.ROLE_ADMIN_MAHASISWA), userController.DeleteNonAdmin)
	}

}
