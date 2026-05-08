package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MataKuliah struct {
	ID   uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	Kode string    `gorm:"not null;type:varchar(12)" json:"kode"`
	Name string    `gorm:"not null;type:varchar(255)" json:"name"`
	Sks  uint      `gorm:"not null;type:int" json:"sks"`
}

func (MataKuliah) TableName() string {
	return "mata_kuliah"
}

func (m *MataKuliah) BeforeCreate(db *gorm.DB) error {
	newId, err := uuid.NewV7()
	if err != nil {
		return err
	}
	m.ID = newId

	return nil
}
