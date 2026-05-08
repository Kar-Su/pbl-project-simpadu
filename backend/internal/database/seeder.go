package database

import (
	"context"
	"web-hosting/internal/database/seeders/seeds"
	"web-hosting/internal/modules/role/repository"

	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	ctx := context.Background()
	roleRepo := repository.NewRoleRepository(db)
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
	roleRepo := repository.NewRoleRepository(db)
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

	return nil
}
