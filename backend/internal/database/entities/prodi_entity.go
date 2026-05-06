package entities

type Prodi struct {
	ID        uint    `gorm:"primaryKey;autoIncrement;type:int" json:"id"`
	Name      string  `gorm:"not null;type:varchar(255);uniqueIndex" json:"name"`
	Jenjang   string  `gorm:"not null;type:enum('D3', 'D4')" json:"jenjang"`
	JurusanID uint    `gorm:"not null;type:int;index" json:"jurusan_id"`
	Jurusan   Jurusan `gorm:"foreignKey:JurusanID" json:"jurusan"`
}

func (Prodi) TableName() string {
	return "prodi"
}
