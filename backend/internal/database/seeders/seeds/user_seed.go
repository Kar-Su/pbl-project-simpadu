package seeds

import (
	"encoding/json"
	"io"
	"os"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/modules/user/dto"

	"gorm.io/gorm"
)

func ListUsersSeed(db *gorm.DB) error {
	jsonFile, err := os.Open("internal/database/seeders/json/users.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)

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

		roleId := dto.RoleNameToRoleID(user.RoleName)
		userEntity.RoleID = roleId
		if user.DetailId != nil {
			userEntity.DetailID = user.DetailId
		}

		if err := db.Where("email = ?", user.Email).FirstOrCreate(&userEntity).Error; err != nil {
			return err
		}
	}

	return nil
}
