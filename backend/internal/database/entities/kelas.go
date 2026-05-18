package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Kelas struct {
	ID              uuid.UUID     `gorm:"primaryKey;type:char(36)" json:"id"`
	Name            string        `gorm:"not null;type:varchar(255)" json:"name"`
	Semester        uint          `gorm:"not null;type:int;default:1" json:"semester"`
	TahunAkademikID uint          `gorm:"type:int" json:"tahun_akademik_id"`
	KurikulumKode   string        `gorm:"type:char(12)" json:"kurikulum_kode"`
	ProdiID         uint          `gorm:"type:int" json:"prodi_id"`
	Kurikulum       Kurikulum     `gorm:"foreignKey:KurikulumKode;references:Kode" json:"kurikulum"`

	TahunAkademik   TahunAkademik `gorm:"foreignKey:TahunAkademikID" json:"tahun_akademik"`
	Prodi           Prodi         `gorm:"foreignKey:ProdiID" json:"prodi"`
	Mahasiswa       []User        `gorm:"many2many:kelas_mahasiswa;foreignKey:ID;joinForeignKey:KelasID;References:DetailID;joinReferences:MahasiswaID" json:"mahasiswa"`
}

func (Kelas) TableName() string {
	return "kelas"
}

func (k *Kelas) BeforeCreate(tx *gorm.DB) error {
	newId, err := uuid.NewV7()
	if err != nil {
		return err
	}
	k.ID = newId
	return nil
}
