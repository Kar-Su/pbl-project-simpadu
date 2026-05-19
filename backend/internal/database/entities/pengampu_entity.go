package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Pengampu struct {
	ID         uuid.UUID  `gorm:"primaryKey;type:char(36)" json:"id"`
	KelasID    uuid.UUID  `gorm:"not null;type:char(36)" json:"kelas_id"`
	MKKode     string     `gorm:"not null;type:char(12)" json:"mk_kode"`
	DosenID    uuid.UUID  `gorm:"not null;type:char(36)" json:"dosen_id"`
	Kelas      Kelas      `gorm:"foreignKey:KelasID" json:"kelas"`
	Dosen      User       `gorm:"foreignKey:DosenID;references:DetailID" json:"dosen"`
	MataKuliah MataKuliah `gorm:"foreignKey:MKKode;references:Kode" json:"mata_kuliah"`
}

func (Pengampu) TableName() string {
	return "pengampu"
}

func (p *Pengampu) BeforeCreate(tx *gorm.DB) error {
	newId, err := uuid.NewV7()
	if err != nil {
		return err
	}

	p.ID = newId
	return nil
}
