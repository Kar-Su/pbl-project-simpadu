package seeds

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/modules/role/dto"

	"gorm.io/gorm"
)

func ListRolesSeed(ctx context.Context, db *gorm.DB) error {
	jsonFile, err := os.Open("internal/database/seeders/json/roles.json")
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
		var roleEntity entities.Role
		roleEntity.Name = role.RoleName
		if err := db.WithContext(ctx).Where(entities.Role{Name: role.RoleName}).FirstOrCreate(&roleEntity).Error; err != nil {
			return err
		}
	}
	return nil
}
