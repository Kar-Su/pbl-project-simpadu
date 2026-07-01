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
	kelasController "web-hosting/internal/modules/kelas/controller"
	kelasRepository "web-hosting/internal/modules/kelas/repository"
	kelasService "web-hosting/internal/modules/kelas/service"
	khsController "web-hosting/internal/modules/khs/controller"
	khsRepo "web-hosting/internal/modules/khs/repository"
	khsService "web-hosting/internal/modules/khs/service"
	kurikulumController "web-hosting/internal/modules/kurikulum/controller"
	kurikulumRepo "web-hosting/internal/modules/kurikulum/repository"
	kurikulumService "web-hosting/internal/modules/kurikulum/service"
	mkController "web-hosting/internal/modules/mk/controller"
	mkRepo "web-hosting/internal/modules/mk/repository"
	mkService "web-hosting/internal/modules/mk/service"
	pengampuController "web-hosting/internal/modules/pengampu/controller"
	pengampuRepo "web-hosting/internal/modules/pengampu/repository"
	pengampuService "web-hosting/internal/modules/pengampu/service"
	presensiController "web-hosting/internal/modules/presensi/controller"
	presensiRepo "web-hosting/internal/modules/presensi/repository"
	presensiService "web-hosting/internal/modules/presensi/service"
	prodiController "web-hosting/internal/modules/prodi/controller"
	prodiRepo "web-hosting/internal/modules/prodi/repository"
	prodiService "web-hosting/internal/modules/prodi/service"
	roleController "web-hosting/internal/modules/role/controller"
	roleRepo "web-hosting/internal/modules/role/repository"
	roleService "web-hosting/internal/modules/role/service"
	userController "web-hosting/internal/modules/user/controller"
	userRepo "web-hosting/internal/modules/user/repository"
	userService "web-hosting/internal/modules/user/service"
	workerController "web-hosting/internal/modules/workers/controller"

	"web-hosting/internal/package/constants"
	"web-hosting/internal/workers"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

func InitDatabases(injector do.Injector) {
	do.ProvideNamed[*gorm.DB](injector, constants.DB, func(i do.Injector) (*gorm.DB, error) {
		return configs.SetUpDatabaseConnection(), nil
	})
}

func InitTestDatabases(injector do.Injector) {
	do.ProvideNamed[*gorm.DB](injector, constants.DB_TEST, func(i do.Injector) (*gorm.DB, error) {
		return configs.SetUpDatabaseTestConnection(), nil
	})
}

func RegisterProviders(injector do.Injector) {
	do.ProvideNamed[authService.JwtService](injector, constants.JWTService, func(i do.Injector) (authService.JwtService, error) {
		return authService.NewJwtService(), nil
	})

	InitDatabases(injector)
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	InitTestDatabases(injector)
	// db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB_TEST)

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
	kelasRepo := kelasRepository.NewKelasRepository(db)
	kelasPivotRepo := kelasRepository.NewKelasMahasiswaRepository(db)
	pengampuRepo := pengampuRepo.NewPengampuRepository(db)
	presensiRepo := presensiRepo.NewPresensiRepository(db)
	khsRepo := khsRepo.NewKhsRepository(db)

	roleService := roleService.NewRoleService(roleRepo, db)
	userService := userService.NewUserService(userRepo, roleService, db)
	// authService := authService.NewAuthService(userRepo, refreshTokenRepo, jwtService, db)
	jurusanService := jurusanService.NewJurusanService(jurusanRepo, db)
	prodiService := prodiService.NewProdiService(prodiRepo, jurusanService, db)
	mkService := mkService.NewMkService(mkRepo, db)
	akademikService := akademikService.NewTahunAkademikService(akademikRepo, db)
	kService := kurikulumService.NewKurikulumService(kRepo, prodiService, db)
	khsService := khsService.NewKhsService(khsRepo, db)
	kPivotService := kurikulumService.NewKurikulumMKService(db, kRepo, kPivotRepo, mkRepo)
	kelasServiceVar := kelasService.NewKelasService(db, kelasRepo, akademikRepo, prodiRepo, kRepo)
	kelasPivotService := kelasService.NewKelasMahasiswaService(db, userRepo, kelasRepo, kelasPivotRepo)
	pengampuService := pengampuService.NewPengampuService(db, pengampuRepo)
	presensiService := presensiService.NewPresensiService(db, presensiRepo)

	do.Provide(injector, func(i do.Injector) (authService.AuthService, error) {
		return authService.NewAuthService(userRepo, refreshTokenRepo, kelasPivotRepo, jwtService, db), nil
	})

	do.Provide(injector, func(i do.Injector) (workers.Schedule, error) {
		authService := do.MustInvoke[authService.AuthService](i)
		return workers.NewSchedule(i, authService, presensiService), nil
	})

	do.Provide(injector, func(i do.Injector) (presensiController.PresensiController, error) {
		return presensiController.NewPresensiController(injector, db, presensiService, userService), nil
	})

	do.Provide(injector, func(i do.Injector) (userController.UserController, error) {
		return userController.NewUserController(i, db, userService, roleService), nil
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
	do.Provide(injector, func(i do.Injector) (kelasController.KelasController, error) {
		return kelasController.NewKelasController(i, db, kelasServiceVar), nil
	})
	do.Provide(injector, func(i do.Injector) (kelasController.KelasMahasiswaController, error) {
		return kelasController.NewKelasMahasiswaController(i, db, kelasPivotService), nil
	})
	do.Provide(injector, func(i do.Injector) (pengampuController.PengampuController, error) {
		return pengampuController.NewPengampuController(i, db, pengampuService), nil
	})
	do.Provide(injector, func(i do.Injector) (workerController.WorkerController, error) {
		return workerController.NewWorkerController(i), nil
	})
	do.Provide(injector, func(i do.Injector) (khsController.KHSController, error) {
		return khsController.NewKHSController(i, db, khsService), nil
	})
}
