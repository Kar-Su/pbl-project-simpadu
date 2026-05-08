package seeds

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"web-hosting/internal/database/entities"

	"gorm.io/gorm"
)

func ListSeedAkademik(ctx context.Context, db *gorm.DB) error {
	jsonFile, err := os.Open("internal/database/seeders/json/tahun-akademik.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var akademikList []entities.TahunAkademik
	if err := json.Unmarshal(jsonData, &akademikList); err != nil {
		return err
	}

	for _, akademik := range akademikList {
		if err := db.WithContext(ctx).Where("id = ?", akademik.ID).FirstOrCreate(&akademik).Error; err != nil {
			return err
		}
	}
	return nil
}
