package database

import (
	"context"
	"web-hosting/internal/database/seeders/seeds"
	pivotKelasRepos "web-hosting/internal/modules/kelas/repository"
	roleRepos "web-hosting/internal/modules/role/repository"

	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	ctx := context.Background()
	roleRepo := roleRepos.NewRoleRepository(db)
	if err := seeds.ListRolesSeed(ctx, db); err != nil {
		return err
	}

	if err := seeds.ListUsersSeed(ctx, db, roleRepo); err != nil {
		return err
	}
	if err := seeds.ListJurusanSeed(ctx, db); err != nil {
		return err
	}

	if err := seeds.ListSeedProdi(ctx, db); err != nil {
		return err
	}

	return nil
}

func SeedDummy(db *gorm.DB) error {
	ctx := context.Background()
	roleRepo := roleRepos.NewRoleRepository(db)
	kelasRepo := pivotKelasRepos.NewKelasRepository(db)
	if err := seeds.ListRolesSeed(ctx, db); err != nil {
		return err
	}

	if err := seeds.ListUsersSeed(ctx, db, roleRepo); err != nil {
		return err
	}

	if err := seeds.ListJurusanSeed(ctx, db); err != nil {
		return err
	}

	if err := seeds.ListSeedProdi(ctx, db); err != nil {
		return err
	}

	if err := seeds.ListMKSeed(ctx, db); err != nil {
		return err
	}

	if err := seeds.ListSeedAkademik(ctx, db); err != nil {
		return err
	}

	if err := seeds.ListKurikulumSeed(ctx, db); err != nil {
		return err
	}

	if err := seeds.ListKurikulumPivotSeed(ctx, db); err != nil {
		return err
	}

	if err := seeds.ListKelasSeed(ctx, db); err != nil {
		return err
	}

	if err := seeds.ListKelasMahasiswaSeed(ctx, db, kelasRepo); err != nil {
		return err
	}

	return nil
}
