package seeds

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/modules/role/repository"
	"web-hosting/internal/modules/user/dto"
	"web-hosting/internal/package/helpers"

	"gorm.io/gorm"
)

func ListUsersSeed(ctx context.Context, db *gorm.DB, roleRepo repository.RoleRepository) error {
	jsonFile, err := os.Open("internal/database/seeders/json/users.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var users []dto.UserAdminCreateRequest
	if err := json.Unmarshal(jsonData, &users); err != nil {
		return err
	}

	for _, user := range users {
		var userEntity entities.User = entities.User{
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
		}

		normRoleName := helpers.NormalizeString(user.RoleName)
		roleId, _ := roleRepo.GetRoleIdByRoleName(ctx, db, normRoleName)
		userEntity.RoleID = roleId
		if user.DetailId != nil {
			userEntity.DetailID = user.DetailId
		}

		if err := db.WithContext(ctx).Where("email = ?", user.Email).FirstOrCreate(&userEntity).Error; err != nil {
			return err
		}
	}

	return nil
}
