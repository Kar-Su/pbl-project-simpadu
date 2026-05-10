package providers

import (
	"web-hosting/internal/configs"
	"web-hosting/internal/database/entities"
	akademikController "web-hosting/internal/modules/akademik/controller"
	akademikRepo "web-hosting/internal/modules/akademik/repository"
	akademikService "web-hosting/internal/modules/akademik/service"
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
	"web-hosting/internal/workers"

	kurikulumController "web-hosting/internal/modules/kurikulum/controller"
	kurikulumRepo "web-hosting/internal/modules/kurikulum/repository"
	kurikulumService "web-hosting/internal/modules/kurikulum/service"

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
	db.SetupJoinTable(&entities.Kelas{}, "Mahasiswa", &entities.KelasMahasiswa{})

	jwtService := do.MustInvokeNamed[authService.JwtService](injector, constants.JWTService)

	userRepo := userRepo.NewUserRepository(db)
	refreshTokenRepo := authRepo.NewRefreshTokenRepository(db)
	roleRepo := roleRepo.NewRoleRepository(db)
	jurusanRepo := jurusanRepo.NewJurusanRepository(db)
	prodiRepo := prodiRepo.NewProdiRepository(db)
	mkRepo := mkRepo.NewMkRepository(db)
	akademikRepo := akademikRepo.NewTahunAkademikRepository(db)
	kRepo := kurikulumRepo.NewKurikulumRepository(db)
	kPivotRepo := kurikulumRepo.NewKurikulumMKRepository(db)

	roleService := roleService.NewRoleService(roleRepo, db)
	userService := userService.NewUserService(userRepo, roleService, db)
	// authService := authService.NewAuthService(userRepo, refreshTokenRepo, jwtService, db)
	jurusanService := jurusanService.NewJurusanService(jurusanRepo, db)
	prodiService := prodiService.NewProdiService(prodiRepo, jurusanService, db)
	mkService := mkService.NewMkService(mkRepo, db)
	akademikService := akademikService.NewTahunAkademikService(akademikRepo, db)
	kService := kurikulumService.NewKurikulumService(kRepo, prodiService, db)
	kPivotService := kurikulumService.NewKurikulumMKService(db, kRepo, kPivotRepo, mkRepo)

	do.Provide(injector, func(i do.Injector) (authService.AuthService, error) {
		return authService.NewAuthService(userRepo, refreshTokenRepo, jwtService, db), nil
	})

	do.Provide(injector, func(i do.Injector) (workers.Schedule, error) {
		authService := do.MustInvoke[authService.AuthService](i)
		return workers.NewSchedule(i, authService), nil
	})

	do.Provide(injector, func(i do.Injector) (userController.UserController, error) {
		return userController.NewUserController(i, userService, roleService), nil
	})

	do.Provide(injector, func(i do.Injector) (authController.AuthController, error) {
		authService := do.MustInvoke[authService.AuthService](i)
		return authController.NewAuthController(i, authService, db), nil
	})

	do.Provide(injector, func(i do.Injector) (roleController.RoleController, error) {
		return roleController.NewRoleController(i, roleService, db), nil
	})

	do.Provide(injector, func(i do.Injector) (kurikulumController.KurikulumController, error) {
		return kurikulumController.NewKurikulumController(i, kService, db), nil
	})
	do.Provide(injector, func(i do.Injector) (kurikulumController.PivotController, error) {
		return kurikulumController.NewPivotController(i, kPivotService, db), nil
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
	do.Provide(injector, func(i do.Injector) (akademikController.TahunAkademikController, error) {
		return akademikController.NewTahunAkademikController(i, akademikService, db), nil
	})
}
