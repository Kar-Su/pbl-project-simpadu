package entities

import "github.com/google/uuid"

type PresensiMahasiswa struct {
	PresensiID  uuid.UUID `gorm:"PrimaryKey;type:char(36);not null"`
	MahasiswaID uuid.UUID `gorm:"PrimaryKey;type:char(36);not null"`
	Status      string    `gorm:"type:enum('hadir', 'sakit', 'izin', 'alpha');default: 'alpha';not null"`
	Mahasiswa   User      `gorm:"foreignKey:MahasiswaID;references:DetailID" json:"mahasiswa"`
}

func (PresensiMahasiswa) TableName() string {
	return "presensi_mahasiswa"
}
