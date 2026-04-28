package seeds

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"web-hosting/internal/modules/role/dto"

	"gorm.io/gorm"
)

func ListRolesSeed(ctx context.Context, db *gorm.DB) error {
	jsonFile, err := os.Open("internal/datavase/seeders/json/roles.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var roles []dto.RoleCreateRequest
	if err := json.Unmarshal(jsonData, &roles); err != nil {
		return err
	}

	for _, role := range roles {
		if err := db.WithContext(ctx).FirstOrCreate(&role).Error; err != nil {
			return err
		}
	}
	return nil
}
