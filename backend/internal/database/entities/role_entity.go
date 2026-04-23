package entities

type Role struct {
	Kode string `gorm:"primaryKey;type:char(4)" json:"kode"`
	Name string `gorm:"not null" json:"name"`
}
