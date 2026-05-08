package seeds

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"web-hosting/internal/database/entities"

	"gorm.io/gorm"
)

func ListKurikulumSeed(ctx context.Context, db *gorm.DB) error {
	jsonFile, err := os.Open("internal/database/seeders/json/kurikulum.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var kurikulumList []entities.Kurikulum
	if err := json.Unmarshal(jsonData, &kurikulumList); err != nil {
		return err
	}

	for _, kurikulum := range kurikulumList {
		if err := db.WithContext(ctx).Where("kode = ?", kurikulum.Kode).FirstOrCreate(&kurikulum).Error; err != nil {
			return err
		}
	}

	return nil
}
