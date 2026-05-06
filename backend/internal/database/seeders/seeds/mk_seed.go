package seeds

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"web-hosting/internal/database/entities"

	"gorm.io/gorm"
)

func ListMKSeed(ctx context.Context, db *gorm.DB) error {
	jsonFile, err := os.Open("internal/database/seeders/json/mk.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var mks []entities.MataKuliah
	if err := json.Unmarshal(jsonData, &mks); err != nil {
		return err
	}

	for _, mk := range mks {
		if err := db.WithContext(ctx).Where("kode = ?", mk.Kode).FirstOrCreate(&mk).Error; err != nil {
			return err
		}
	}

	return nil
}
