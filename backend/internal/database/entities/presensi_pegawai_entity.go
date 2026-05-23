package entities

import "github.com/google/uuid"

type PresensiPegawai struct {
	PresensiID uuid.UUID `gorm:"PrimaryKey;type:char(36);not null"`
	PegawaiID  uuid.UUID `gorm:"PrimaryKey;type:char(36);not null"`
	Status     string    `gorm:"type:enum('hadir', 'sakit', 'izin', 'alpha');default: 'alpha';not null"`
	Pegawai    User      `gorm:"foreignKey:PegawaiID;references:DetailID" json:"pegawai"`
}

func (PresensiPegawai) TableName() string {
	return "presensi_pegawai"
}
