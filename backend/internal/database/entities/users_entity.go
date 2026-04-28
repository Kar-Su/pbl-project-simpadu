package entities

import (
	"web-hosting/internal/package/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	Name     string    `gorm:"type:varchar(255);not null" json:"name"`
	Email    string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	Password string    `gorm:"type:varchar(255);not null" json:"password"`
	ImageUrl *string   `gorm:"type:varchar(255)" json:"image_url"`
	RoleID   uint      `gorm:"type:int;not null" json:"role_id"`
	Role     Role      `gorm:"foreignKey:RoleID;references:ID" json:"role"`
	DetailID *uint     `gorm:"type:int;index" json:"detail_id"`

	Timestamp
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	newId, err := uuid.NewV7()
	if err != nil {
		return err
	}
	u.ID = newId
	if u.Password != "" {
		u.Password, err = helpers.HashPassword(u.Password)
		if err != nil {
			return err
		}
	}

	return
}
