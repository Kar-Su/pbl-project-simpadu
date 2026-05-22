package entities

import (
	"time"
	"web-hosting/internal/package/types"

	"github.com/google/uuid"
)

type Presensi struct {
	ID                uuid.UUID           `gorm:"primaryKey;type:char(36)"`
	Tipe              string              `gorm:"type:enum('mahasiswa', 'pegawai');not null"`
	PengampuID        uuid.UUID           `gorm:"type:char(36);not null"`
	CreatedAt         types.DateOnly      `gorm:"type:date;default:CURRENT_DATE()"`
	UpdatedAt         time.Time           `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	PresensiMahasiswa []PresensiMahasiswa `gorm:"foreignKey:PresensiID"`
	PresensiPegawai   []PresensiPegawai   `gorm:"foreignKey:PresensiID"`
}

func (Presensi) TableName() string {
	return "presensi"
}
