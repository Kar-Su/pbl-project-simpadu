package repository

import (
	"context"
	"errors"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/modules/khs/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KhsRepository interface {
	Create(ctx context.Context, tx *gorm.DB, mahasiswaID uuid.UUID, nilai *entities.NilaiMk) error
	GetKHS(ctx context.Context, tx *gorm.DB, filter *dto.FilterQuery, offset, limit int) ([]entities.Khs, error)
}

type khsRepository struct {
	db *gorm.DB
}

func NewKhsRepository(db *gorm.DB) KhsRepository {
	return &khsRepository{db: db}
}

func (r *khsRepository) Create(ctx context.Context, tx *gorm.DB, mahasiswaID uuid.UUID, nilai *entities.NilaiMk) error {
	if tx == nil {
		tx = r.db.WithContext(ctx)
	} else {
		tx = tx.WithContext(ctx)
	}

	return tx.Transaction(func(dbTx *gorm.DB) error {
		var result struct {
			Semester uint
		}
		err := dbTx.Table("pengampu").
			Select("kelas.semester").
			Joins("JOIN kelas ON kelas.id = pengampu.kelas_id").
			Where("pengampu.id = ?", nilai.PengampuID).
			Scan(&result).Error

		if err != nil {
			return err
		}
		if result.Semester == 0 {
			return errors.New("pengampu or kelas not found")
		}

		khs := entities.Khs{
			Semester:    result.Semester,
			MahasiswaID: mahasiswaID,
		}
		if err := dbTx.Where(entities.Khs{Semester: result.Semester, MahasiswaID: mahasiswaID}).
			FirstOrCreate(&khs).Error; err != nil {
			return err
		}

		switch {
		case nilai.TotalNilai >= 80:
			nilai.GradeNilai = "A"
		case nilai.TotalNilai >= 70:
			nilai.GradeNilai = "B"
		case nilai.TotalNilai >= 60:
			nilai.GradeNilai = "C"
		case nilai.TotalNilai >= 50:
			nilai.GradeNilai = "D"
		default:
			nilai.GradeNilai = "E"
		}

		nilai.KhsID = khs.ID
		if err := dbTx.Create(nilai).Error; err != nil {
			return err
		}

		var ipsResult struct {
			TotalSks  uint    `gorm:"column:total_sks"`
			TotalMutu float32 `gorm:"column:total_mutu"`
		}
		if err := dbTx.Table("nilai_mk").
			Select(`
				COALESCE(SUM(mata_kuliah.sks), 0) as total_sks,
				COALESCE(SUM(
					CASE nilai_mk.grade_nilai
						WHEN 'A' THEN 4.0 WHEN 'B' THEN 3.0 WHEN 'C' THEN 2.0 WHEN 'D' THEN 1.0 ELSE 0.0
					END * mata_kuliah.sks
				), 0) as total_mutu
			`).
			Joins("JOIN pengampu ON pengampu.id = nilai_mk.pengampu_id").
			Joins("JOIN mata_kuliah ON mata_kuliah.kode = pengampu.mk_kode").
			Where("nilai_mk.khs_id = ?", khs.ID).
			Scan(&ipsResult).Error; err != nil {
			return err
		}

		var ips float32 = 0
		if ipsResult.TotalSks > 0 {
			ips = ipsResult.TotalMutu / float32(ipsResult.TotalSks)
		}

		var ipkResult struct {
			TotalSks  uint    `gorm:"column:total_sks"`
			TotalMutu float32 `gorm:"column:total_mutu"`
		}
		if err := dbTx.Table("nilai_mk").
			Select(`
				COALESCE(SUM(mata_kuliah.sks), 0) as total_sks,
				COALESCE(SUM(
					CASE nilai_mk.grade_nilai
						WHEN 'A' THEN 4.0 WHEN 'B' THEN 3.0 WHEN 'C' THEN 2.0 WHEN 'D' THEN 1.0 ELSE 0.0
					END * mata_kuliah.sks
				), 0) as total_mutu
			`).
			Joins("JOIN khs ON khs.id = nilai_mk.khs_id").
			Joins("JOIN pengampu ON pengampu.id = nilai_mk.pengampu_id").
			Joins("JOIN mata_kuliah ON mata_kuliah.kode = pengampu.mk_kode").
			Where("khs.mahasiswa_id = ?", mahasiswaID).
			Scan(&ipkResult).Error; err != nil {
			return err
		}

		var ipk float32 = 0
		if ipkResult.TotalSks > 0 {
			ipk = ipkResult.TotalMutu / float32(ipkResult.TotalSks)
		}

		if err := dbTx.Model(&entities.Khs{}).
			Where("id = ?", khs.ID).
			Updates(map[string]interface{}{
				"ips": ips,
				"ipk": ipk,
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *khsRepository) GetKHS(ctx context.Context, tx *gorm.DB, filter *dto.FilterQuery, offset, limit int) ([]entities.Khs, error) {
	db := tx
	if db == nil {
		db = r.db.WithContext(ctx)
	} else {
		db = db.WithContext(ctx)
	}

	var listKhs []entities.Khs
	query := db.Model(&entities.Khs{})

	if filter != nil {
		if filter.Semester != 0 {
			query = query.Where("khs.semester = ?", filter.Semester)
		}

		if filter.ProdiName != "" {
			query = query.
				Joins("JOIN nilai_mk ON nilai_mk.khs_id = khs.id").
				Joins("JOIN pengampu ON pengampu.id = nilai_mk.pengampu_id").
				Joins("JOIN kelas ON kelas.id = pengampu.kelas_id").
				Joins("JOIN prodi ON prodi.id = kelas.prodi_id").
				Where("LOWER(prodi.name) LIKE LOWER(?)", "%"+filter.ProdiName+"%").
				Distinct()
		}
	}

	query = query.
		Preload("Mahasiswa").
		Preload("NilaiMk").
		Preload("NilaiMk.Pengampu.MataKuliah").
		Preload("NilaiMk.Pengampu.Kelas.Prodi").
		Preload("NilaiMk.Pengampu.Kelas.TahunAkademik").
		Preload("NilaiMk.Pengampu.Kelas.Kurikulum")

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&listKhs).Error; err != nil {
		return nil, err
	}

	return listKhs, nil
}
