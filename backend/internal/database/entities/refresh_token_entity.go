package entities

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"primaryKey;type:binary(16)" json:"id"`
	UserID    uuid.UUID `gorm:"type:binary(16);not null;index" json:"user_id"`
	Token     string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"token"`
	ExpiredAt time.Time `gorm:"type:timestamp;not null" json:"expired_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`

	Timestamp
}

func (r *RefreshToken) BeforeCreate() (err error) {
	r.ID, err = uuid.NewV7()
	if err != nil {
		return err
	}

	return
}
