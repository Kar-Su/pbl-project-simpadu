package seeds

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"web-hosting/internal/database/entities"

	"gorm.io/gorm"
)

func ListKurikulumPivotSeed(ctx context.Context, db *gorm.DB) error {
	jsonFile, err := os.Open("internal/database/seeders/json/pivot_kurikulum_mk.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var pivotKurikulumMk []entities.KurikulumMK
	if err := json.Unmarshal(jsonData, &pivotKurikulumMk); err != nil {
		return err
	}

	for _, pivot := range pivotKurikulumMk {
		if err := db.WithContext(ctx).Where("kurikulum_kode = ? AND mk_kode = ?", pivot.KurikulumKode, pivot.MkKode).FirstOrCreate(&pivot).Error; err != nil {
			return err
		}
	}

	return nil
}
