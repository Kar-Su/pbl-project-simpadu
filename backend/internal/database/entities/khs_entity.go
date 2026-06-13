package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Khs struct {
	ID          uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	Semester    uint      `gorm:"type:int;not null" json:"semester"`
	IPS         float32   `gorm:"type:decimal(3,2);default:0;column:ips" json:"ips"`
	IPK         float32   `gorm:"type:decimal(3,2);default:0;column:ipk" json:"ipk"`
	MahasiswaID uuid.UUID `gorm:"type:char(36);not null" json:"mahasiswa_id"`
	Mahasiswa   User      `gorm:"foreignKey:MahasiswaId;references:DetailID" json:"mahasiswa"`
	NilaiMk     []NilaiMk `gorm:"foreignKey:KhsID" json:"nilai,omitempty"`
}

func (Khs) TableName() string {
	return "khs"
}

func (khs *Khs) BeforeCreate(tx *gorm.DB) error {
	if khs.ID == uuid.Nil {
		newID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		khs.ID = newID
	}
	return nil
}

type NilaiMk struct {
	ID         uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	TotalNilai float32   `gorm:"type:decimal(5,2);default:0;check:total_nilai >= 0 AND total_nilai <= 100" json:"total_nilai"`
	GradeNilai string    `gorm:"type:enum('A', 'B', 'C', 'D', 'E')" json:"grade_nilai"`
	KhsID      uuid.UUID `gorm:"type:char(36);not null;uniqueIndex:idx_khs_pengampu" json:"khs_id"`
	PengampuID uuid.UUID `gorm:"type:char(36);not null;uniqueIndex:idx_khs_pengampu" json:"pengampu_id"`
	Pengampu   Pengampu  `gorm:"foreignKey:PengampuID" json:"pengampu"`
}

func (NilaiMk) TableName() string {
	return "nilai_mk"
}

func (n *NilaiMk) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		newID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		n.ID = newID
	}
	return nil
}
