package seeds

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"web-hosting/internal/database/entities"

	"gorm.io/gorm"
)

func ListPengampuSeed(ctx context.Context, db *gorm.DB) error {
	jsonFile, err := os.Open("internal/database/seeders/json/pengampu.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var pengampu []entities.Pengampu
	if err := json.Unmarshal(jsonData, &pengampu); err != nil {
		return err
	}

	for _, p := range pengampu {
		if err := db.WithContext(ctx).Create(&p).Error; err != nil {
			return err
		}
	}

	return nil
}
