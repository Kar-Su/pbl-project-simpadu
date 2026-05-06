package seeds

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"web-hosting/internal/database/entities"

	"gorm.io/gorm"
)

func ListSeedProdi(ctx context.Context, db *gorm.DB) error {
	jsonFile, err := os.Open("internal/database/seeders/json/prodi.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var prodis []entities.Prodi
	if err := json.Unmarshal(jsonData, &prodis); err != nil {
		return err
	}

	for _, prodi := range prodis {
		if err := db.WithContext(ctx).Where("name = ?", prodi.Name).FirstOrCreate(&prodi).Error; err != nil {
			return err
		}
	}

	return nil
}
