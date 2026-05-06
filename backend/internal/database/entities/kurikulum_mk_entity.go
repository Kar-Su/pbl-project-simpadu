package entities

type KurikulumMK struct {
	KurikulumKode string `gorm:"primaryKey;type:char(12);"`
	MKKode        string `gorm:"primaryKey;type:char(12);"`
	Semester      int    `gorm:"type:int;not null;"`
	Wajib         bool   `gorm:"type:bool;not null;default:false"`
}

func (KurikulumMK) TableName() string {
	return "kurikulum_mk"
}
