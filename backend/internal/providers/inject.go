package providers

import (
	"web-hosting/internal/configs"
	"web-hosting/internal/package/constants"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

func InitDatabases(injector do.Injector) {
	do.ProvideNamed[*gorm.DB](injector, constants.DB, func(i do.Injector) (*gorm.DB, error) {
		return configs.SetUpDatabaseConnection(), nil
	})
}

func RegisterProviders(injector do.Injector) {
	InitDatabases(injector)

	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	_ = db
}
