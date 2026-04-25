package script

import (
	"log"
	"os"
	"web-hosting/internal/database"
	"web-hosting/internal/package/constants"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

func Commands(injector do.Injector) bool {
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	seed := false

	if arg := os.Args[1]; arg == "--seed" {
		seed = true
	}

	if seed {
		if err := database.Seeder(db); err != nil {
			log.Fatalf("error migration seeder: %v", err)
		}
		log.Println("seeder completed successfully")
	}

	return false
}
