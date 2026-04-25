package database

import (
	"web-hosting/internal/database/seeders/seeds"

	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	if err := seeds.ListUsersSeed(db); err != nil {
		return err
	}

	return nil
}
