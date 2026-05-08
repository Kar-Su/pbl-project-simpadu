package entities

import karTypes "web-hosting/internal/package/types"

type TahunAkademik struct {
	ID           uint              `gorm:"primaryKey;type:int" json:"id"`
	TipeSemester string            `gorm:"not null;type:enum('ganjil', 'genap')" json:"tipe_semester"`
	TahunAwal    karTypes.DateOnly `gorm:"not null;type:date;uniqueIndex" json:"tahun_awal"`
	TahunAkhir   karTypes.DateOnly `gorm:"not null;type:date;uniqueIndex" json:"tahun_akhir"`
	Status       string            `gorm:"not null;type:enum('aktif', 'nonaktif');default:'aktif'" json:"status"`
}

func (TahunAkademik) TableName() string {
	return "tahun_akademik"
}
