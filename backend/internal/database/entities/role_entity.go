package entities

type Role struct {
	ID   uint   `gorm:"primaryKey;autoIncrement;type:int" json:"id"`
	Name string `gorm:"not null" json:"name"`
}
