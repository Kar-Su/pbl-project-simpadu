package seeds

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"web-hosting/internal/database/entities"

	"gorm.io/gorm"
)

func ListKelasSeed(ctx context.Context, db *gorm.DB) error {
	jsonFile, err := os.Open("internal/database/seeders/json/kelas.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var kelasList []entities.Kelas
	if err := json.Unmarshal(jsonData, &kelasList); err != nil {
		return err
	}

	for _, kelas := range kelasList {
		if err := db.WithContext(ctx).Create(&kelas).Error; err != nil {
			return err
		}
	}

	return nil
}
