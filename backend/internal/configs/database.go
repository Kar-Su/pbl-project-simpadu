package configs

import (
	"fmt"
	"web-hosting/internal/package/env"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetUpDatabaseConnection() *gorm.DB {
	dbUser := env.GetWithDefault[string]("DB_USERNAME", "root")
	dbPass := env.GetWithDefault[string]("DB_PASSWORD", "root")
	dbHost := env.GetWithDefault[string]("DB_HOST", "db")
	dbName := env.GetWithDefault[string]("GO_DB", "main")
	dbPort := env.GetWithDefault[int]("DB_PORT", 3306)

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v", dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: SetupLogger(),
	})
	if err != nil {
		panic(err)
	}

	return db
}
func SetUpDatabaseTestConnection() *gorm.DB {
	dbUser := env.GetWithDefault[string]("DB_USERNAME", "root")
	dbPass := env.GetWithDefault[string]("DB_PASSWORD", "root")
	dbHost := env.GetWithDefault[string]("DB_HOST", "db")
	dbName := env.GetWithDefault[string]("TEST_DB", "test")
	dbPort := env.GetWithDefault[int]("DB_PORT", 3306)

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v", dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: SetupLogger(),
	})
	if err != nil {
		panic(err)
	}

	return db
}
