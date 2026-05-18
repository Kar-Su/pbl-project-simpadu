package seeds

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/modules/kelas/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ListKelasMahasiswaSeed(ctx context.Context, db *gorm.DB, pivotRepo repository.KelasRepository) error {
	jsonFile, err := os.Open("internal/database/seeders/json/kelas_mahasiswa.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var seedRequests []KelasMahasiswaSeedRequest
	if err := json.Unmarshal(jsonData, &seedRequests); err != nil {
		return err
	}

	for _, request := range seedRequests {
		var kelasIdString string
		if err := db.WithContext(ctx).Table("kelas").Select("id").Where("name = ?", request.KelasName).Limit(1).Scan(&kelasIdString).Error; err != nil {
			return err
		}

		entityPivot := entities.KelasMahasiswa{
			KelasID:     uuid.MustParse(kelasIdString),
			MahasiswaID: request.MahasiswaID,
		}

		if err := db.WithContext(ctx).Create(&entityPivot).Error; err != nil {
			return err
		}

	}

	return nil
}
