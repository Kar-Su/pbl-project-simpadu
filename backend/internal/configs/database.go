package configs

import (
	"fmt"
	"log"
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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: SetupLogger(),
	})

	if err != nil {
		log.Printf("Error Init Database: %s", dsn)
		panic(err)
	}

	return db
}

func SetUpDatabaseTestConnection() *gorm.DB {
	// dbUser := env.GetWithDefault[string]("DB_USERNAME", "root")
	dbUser := "root"
	// dbPass := env.GetWithDefault[string]("DB_PASSWORD", "root")
	dbPass := env.GetWithDefault[string]("DB_ROOT_PASSWORD", "root")
	dbHost := env.GetWithDefault[string]("DB_HOST", "db")
	dbName := env.GetWithDefault[string]("TEST_DB", "test")
	dbPort := env.GetWithDefault[int]("DB_PORT", 3306)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: SetupLogger(),
	})

	if err != nil {
		log.Printf("Error Init Database: %s", dsn)
		panic(err)
	}

	return db
}
