package entities

type Jurusan struct {
	ID   uint   `gorm:"primaryKey;autoIncrement;type:int" json:"id"`
	Name string `gorm:"type:varchar(255);not null;uniqueIndex" json:"name"`
}

func (Jurusan) TableName() string {
	return "jurusan"
}
