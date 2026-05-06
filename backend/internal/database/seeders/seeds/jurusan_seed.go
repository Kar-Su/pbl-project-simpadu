package seeds

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"web-hosting/internal/database/entities"

	"gorm.io/gorm"
)

func ListJurusanSeed(ctx context.Context, db *gorm.DB) error {
	jsonFile, err := os.Open("internal/database/seeders/json/jurusan.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var jurusanList []entities.Jurusan
	if err := json.Unmarshal(jsonData, &jurusanList); err != nil {
		return err
	}

	for _, jurusan := range jurusanList {
		if err := db.WithContext(ctx).Where("name = ?", jurusan.Name).FirstOrCreate(&jurusan).Error; err != nil {
			return err
		}
	}

	return nil
}
