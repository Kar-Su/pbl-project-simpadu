package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Kurikulum struct {
	ID         uuid.UUID    `gorm:"primaryKey;type:char(36)" json:"id"`
	Kode       string       `gorm:"not null;type:char(12);uniqueIndex" json:"kode"`
	Name       string       `gorm:"not null;type:varchar(255)" json:"nama"`
	ProdiID    uint         `gorm:"not null;type:int" json:"prodi_id"`
	Prodi      Prodi        `gorm:"foreignKey:ProdiID" json:"prodi"`
	MataKuliah []MataKuliah `gorm:"many2many:kurikulum_mk" json:"mata_kuliah"`
}

func (Kurikulum) TableName() string {
	return "kurikulum"
}

func (k *Kurikulum) BeforeCreate(db *gorm.DB) error {
	newId, err := uuid.NewV7()
	if err != nil {
		return err
	}
	k.ID = newId

	return nil
}
