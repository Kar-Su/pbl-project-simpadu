package providers

import (
	"web-hosting/internal/configs"
	"web-hosting/internal/database/entities"
	authController "web-hosting/internal/modules/auth/controller"
	authRepo "web-hosting/internal/modules/auth/repository"
	authService "web-hosting/internal/modules/auth/service"
	jurusanController "web-hosting/internal/modules/jurusan/controller"
	jurusanRepo "web-hosting/internal/modules/jurusan/repository"
	jurusanService "web-hosting/internal/modules/jurusan/service"
	mkController "web-hosting/internal/modules/mk/controller"
	mkRepo "web-hosting/internal/modules/mk/repository"
	mkService "web-hosting/internal/modules/mk/service"
	prodiController "web-hosting/internal/modules/prodi/controller"
	prodiRepo "web-hosting/internal/modules/prodi/repository"
	prodiService "web-hosting/internal/modules/prodi/service"
	roleController "web-hosting/internal/modules/role/controller"
	roleRepo "web-hosting/internal/modules/role/repository"
	roleService "web-hosting/internal/modules/role/service"
	userController "web-hosting/internal/modules/user/controller"
	userRepo "web-hosting/internal/modules/user/repository"
	userService "web-hosting/internal/modules/user/service"

	"web-hosting/internal/package/constants"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

func InitDatabases(injector do.Injector) {
	do.ProvideNamed[*gorm.DB](injector, constants.DB, func(i do.Injector) (*gorm.DB, error) {
		return configs.SetUpDatabaseConnection(), nil
	})
}

func InitTestDatabases(injector do.Injector) {
	do.ProvideNamed[*gorm.DB](injector, "db_test", func(i do.Injector) (*gorm.DB, error) {
		return configs.SetUpDatabaseTestConnection(), nil
	})
}

func RegisterProviders(injector do.Injector) {
	// InitDatabases(injector)
	InitTestDatabases(injector)
	do.ProvideNamed[authService.JwtService](injector, constants.JWTService, func(i do.Injector) (authService.JwtService, error) {
		return authService.NewJwtService(), nil
	})

	db := do.MustInvokeNamed[*gorm.DB](injector, "db_test")
	db.SetupJoinTable(&entities.Kurikulum{}, "MataKuliah", &entities.KurikulumMK{})

	jwtService := do.MustInvokeNamed[authService.JwtService](injector, constants.JWTService)

	userRepo := userRepo.NewUserRepository(db)
	refreshTokenRepo := authRepo.NewRefreshTokenRepository(db)
	roleRepo := roleRepo.NewRoleRepository(db)
	jurusanRepo := jurusanRepo.NewJurusanRepository(db)
	prodiRepo := prodiRepo.NewProdiRepository(db)
	mkRepo := mkRepo.NewMkRepository(db)

	roleService := roleService.NewRoleService(roleRepo, db)
	userService := userService.NewUserService(userRepo, roleService, db)
	authService := authService.NewAuthService(userRepo, refreshTokenRepo, jwtService, db)
	jurusanService := jurusanService.NewJurusanService(jurusanRepo, db)
	prodiService := prodiService.NewProdiService(prodiRepo, jurusanService, db)
	mkService := mkService.NewMkService(mkRepo, db)

	do.Provide(injector, func(i do.Injector) (userController.UserController, error) {
		return userController.NewUserController(i, userService, roleService), nil
	})

	do.Provide(injector, func(i do.Injector) (authController.AuthController, error) {
		return authController.NewAuthController(i, authService, db), nil
	})

	do.Provide(injector, func(i do.Injector) (roleController.RoleController, error) {
		return roleController.NewRoleController(i, roleService, db), nil
	})

	do.Provide(injector, func(i do.Injector) (jurusanController.JurusanController, error) {
		return jurusanController.NewJurusanController(i, jurusanService, db), nil
	})

	do.Provide(injector, func(i do.Injector) (prodiController.ProdiController, error) {
		return prodiController.NewProdiController(i, prodiService, db), nil
	})
	do.Provide(injector, func(i do.Injector) (mkController.MkController, error) {
		return mkController.NewMkController(i, mkService, db), nil
	})
}
