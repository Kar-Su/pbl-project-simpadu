package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID    uuid.UUID `gorm:"type:char(36);not null;index" json:"user_id"`
	Token     string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"token"`
	ExpiredAt time.Time `gorm:"type:timestamp;not null" json:"expired_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`

	Timestamp
}

func (r *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	newId, err := uuid.NewV7()
	if err != nil {
		return err
	}

	r.ID = newId
	return
}
