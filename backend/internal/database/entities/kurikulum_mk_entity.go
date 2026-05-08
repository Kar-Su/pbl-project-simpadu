package entities

type KurikulumMK struct {
	KurikulumKode string     `gorm:"primaryKey;type:char(12);" json:"kurikulum_kode"`
	MkKode        string     `gorm:"primaryKey;type:char(12);" json:"mk_kode"`
	Semester      int        `gorm:"type:int;not null;" json:"semester"`
	Wajib         bool       `gorm:"type:bool;not null;default:false" json:"wajib"`
	MataKuliah    MataKuliah `gorm:"foreignKey:MkKode;references:Kode" json:"mata_kuliah"`
}

func (KurikulumMK) TableName() string {
	return "kurikulum_mk"
}
