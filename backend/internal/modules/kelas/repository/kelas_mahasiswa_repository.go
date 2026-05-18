package repository

import (
	"context"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/package/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	KelasMahasiswaRepository interface {
		Create(ctx context.Context, tx *gorm.DB, entity entities.KelasMahasiswa) error
		Delete(ctx context.Context, tx *gorm.DB, mahasiswaId uuid.UUID, kelasId uuid.UUID) error
		GetAllKelasMahasiswa(ctx context.Context, tx *gorm.DB, mahasiswaId uuid.UUID) ([]entities.KelasMahasiswa, error)
		GetMahasiswaByKelasId(ctx context.Context, tx *gorm.DB, kelasId uuid.UUID) ([]entities.KelasMahasiswa, error)
		CheckMahasiswaAlreadyAssigned(ctx context.Context, tx *gorm.DB, mahasiswaId uuid.UUID) (bool, error)
		GetKelasIdByMahasiswa(ctx context.Context, tx *gorm.DB, mahasiswaId *uuid.UUID) (uuid.UUID, error)
	}

	kelasMahasiswaRepository struct {
		db *gorm.DB
	}
)

func NewKelasMahasiswaRepository(db *gorm.DB) KelasMahasiswaRepository {
	return &kelasMahasiswaRepository{db: db}
}

func (k *kelasMahasiswaRepository) Create(ctx context.Context, tx *gorm.DB, entity entities.KelasMahasiswa) error {
	if tx == nil {
		tx = k.db
	}

	if err := tx.WithContext(ctx).Create(&entity).Error; err != nil {
		return err
	}
	return nil
}

func (k *kelasMahasiswaRepository) Delete(ctx context.Context, tx *gorm.DB, mahasiswaId uuid.UUID, kelasId uuid.UUID) error {
	if tx == nil {
		tx = k.db
	}

	result := tx.WithContext(ctx).
		Where("mahasiswa_id = ? AND kelas_id = ?", mahasiswaId, kelasId).
		Delete(&entities.KelasMahasiswa{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (k *kelasMahasiswaRepository) GetMahasiswaByKelasId(ctx context.Context, tx *gorm.DB, kelasId uuid.UUID) ([]entities.KelasMahasiswa, error) {
	if tx == nil {
		tx = k.db
	}

	var result []entities.KelasMahasiswa
	if err := tx.WithContext(ctx).
		Preload("Mahasiswa", helpers.SelectFields("id, detail_id, email, name")).
		Where("kelas_id = ?", kelasId).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (k *kelasMahasiswaRepository) GetAllKelasMahasiswa(ctx context.Context, tx *gorm.DB, mahasiswaId uuid.UUID) ([]entities.KelasMahasiswa, error) {
	if tx == nil {
		tx = k.db
	}

	var result []entities.KelasMahasiswa
	if err := tx.WithContext(ctx).
		Preload("Mahasiswa", helpers.SelectFields("id, detail_id, email, name")).
		Preload("Kelas", helpers.SelectFields("id, name, semester")).
		Where("mahasiswa_id = ?", mahasiswaId).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (k *kelasMahasiswaRepository) CheckMahasiswaAlreadyAssigned(ctx context.Context, tx *gorm.DB, mahasiswaId uuid.UUID) (bool, error) {
	if tx == nil {
		tx = k.db
	}

	var count int64
	if err := tx.WithContext(ctx).
		Model(&entities.KelasMahasiswa{}).
		Where("mahasiswa_id = ?", mahasiswaId).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (k *kelasMahasiswaRepository) GetKelasIdByMahasiswa(ctx context.Context, tx *gorm.DB, mahasiswaId *uuid.UUID) (uuid.UUID, error) {
	if mahasiswaId == nil {
		return uuid.Nil, nil
	}
	if tx == nil {
		tx = k.db
	}

	var result entities.KelasMahasiswa
	if err := tx.WithContext(ctx).
		Where("mahasiswa_id = ?", *mahasiswaId).
		First(&result).Error; err != nil {
		return uuid.Nil, err
	}

	return result.KelasID, nil
}
