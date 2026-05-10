package repository

import (
	"context"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/package/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	KelasRepository interface {
		Create(ctx context.Context, tx *gorm.DB, entity entities.Kelas) error
		Update(ctx context.Context, tx *gorm.DB, id uuid.UUID, entity entities.Kelas) error
		Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
		GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entities.Kelas, error)
		GetByProdiID(ctx context.Context, tx *gorm.DB, prodiID uint) ([]entities.Kelas, error)
		GetNonPreload(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entities.Kelas, error)
		CheckKelasById(ctx context.Context, tx *gorm.DB, id uuid.UUID) (bool, error)
	}
	kelasRepository struct {
		db *gorm.DB
	}
)

func NewKelasRepository(db *gorm.DB) KelasRepository {
	return &kelasRepository{db: db}
}

func (r *kelasRepository) Create(ctx context.Context, tx *gorm.DB, entity entities.Kelas) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&entity).Error; err != nil {
		return err
	}
	return nil
}

func (r *kelasRepository) Update(ctx context.Context, tx *gorm.DB, id uuid.UUID, entity entities.Kelas) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Where("id = ?", id).Updates(&entity).Error; err != nil {
		return err
	}
	return nil
}

func (r *kelasRepository) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	result := tx.WithContext(ctx).Where("id = ?", id).Delete(&entities.Kelas{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *kelasRepository) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entities.Kelas, error) {
	if tx == nil {
		tx = r.db
	}

	var kelas entities.Kelas
	if err := tx.WithContext(ctx).
		Preload("Prodi.Jurusan").
		Preload("Kurikulum", helpers.SelectFields("kode, name")).
		Preload("KurikulumMK", "semester = (SELECT semester FROM kelas WHERE id = ?)", id).
		Preload("KurikulumMK.MataKuliah", helpers.SelectFields("kode", "name", "sks")).
		Preload("TahunAkademik").
		Preload("Mahasiswa", helpers.SelectFields("detail_id, name, email")).
		Where("id = ?", id).
		First(&kelas).Error; err != nil {
		return entities.Kelas{}, err
	}

	return kelas, nil
}

func (r *kelasRepository) GetByProdiID(ctx context.Context, tx *gorm.DB, prodiID uint) ([]entities.Kelas, error) {
	if tx == nil {
		tx = r.db
	}

	var kelas []entities.Kelas
	if err := tx.WithContext(ctx).
		Preload("Kurikulum", helpers.SelectFields("kode", "name")).
		Preload("KurikulumMK", func(db *gorm.DB) *gorm.DB {
			return db.Joins("JOIN kelas ON kelas.kurikulum_kode = kurikulum_mk.kurikulum_kode AND kelas.semester = kurikulum_mk.semester")
		}).
		Preload("KurikulumMK.MataKuliah", helpers.SelectFields("kode", "name", "sks")).
		Preload("TahunAkademik").
		Preload("Mahasiswa", helpers.SelectFields("detail_id", "name", "email")).
		Where("prodi_id = ?", prodiID).
		Find(&kelas).Error; err != nil {
		return nil, err
	}

	return kelas, nil
}

func (r *kelasRepository) GetNonPreload(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entities.Kelas, error) {
	if tx == nil {
		tx = r.db
	}

	var kelas entities.Kelas
	if err := tx.WithContext(ctx).
		Where("id = ?", id).
		First(&kelas).Error; err != nil {
		return entities.Kelas{}, err
	}

	return kelas, nil
}

func (r *kelasRepository) CheckKelasById(ctx context.Context, tx *gorm.DB, id uuid.UUID) (bool, error) {
	if tx == nil {
		tx = r.db
	}

	var count int64
	if err := tx.WithContext(ctx).
		Where("id = ?", id).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
