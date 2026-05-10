package entities

import "github.com/google/uuid"

type KelasMahasiswa struct {
	KelasID     uuid.UUID `gorm:"primaryKey;type:char(36)" json:"kelas_id"`
	MahasiswaID uuid.UUID `gorm:"primaryKey;type:char(36)" json:"mahasiswa_id"`
	Mahasiswa   User      `gorm:"foreignKey:MahasiswaID" json:"mahasiswa"`
	Kelas       Kelas     `gorm:"foreignKey:KelasID" json:"kelas"`
}

func (KelasMahasiswa) TableName() string {
	return "kelas_mahasiswa"
}
