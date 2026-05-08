package helpers

import "gorm.io/gorm"

func SelectFields(columns ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(columns)
	}
}
