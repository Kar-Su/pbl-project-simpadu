package script

import (
	"log"
	"os"
	"slices"
	"web-hosting/internal/database"
	"web-hosting/internal/package/constants"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

func Commands(injector do.Injector) bool {
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB_TEST)
	seed := false
	dummy := false

	if arg := os.Args[1]; arg == "--seed" {
		seed = true
	}

	if slices.Contains(os.Args, "--dummy") {
		dummy = true
	}

	if seed {
		if err := database.Seeder(db); err != nil {
			log.Fatalf("error migration seeder: %v", err)
		}
		log.Println("seeder completed successfully")
	}

	if dummy {
		if err := database.SeedDummy(db); err != nil {
			log.Fatalf("error seeding dummy: %v", err)
		}
		log.Println("dummy seed completed successfully")
	}

	return false
}
