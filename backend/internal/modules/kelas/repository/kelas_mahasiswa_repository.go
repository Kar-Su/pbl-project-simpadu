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
		return nil
	}

	if err := tx.WithContext(ctx).Create(&entity).Error; err != nil {
		return err
	}
	return nil
}

func (k *kelasMahasiswaRepository) Delete(ctx context.Context, tx *gorm.DB, mahasiswaId uuid.UUID, kelasId uuid.UUID) error {
	if tx == nil {
		return nil
	}

	result := tx.WithContext(ctx).Delete(&entities.KelasMahasiswa{}, mahasiswaId)
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
		return nil, nil
	}

	var entities []entities.KelasMahasiswa
	if err := tx.WithContext(ctx).
		Preload("Mahasiswa", helpers.SelectFields("detail_id, email, name")).
		Where("kelas_id = ?", kelasId).Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (k *kelasMahasiswaRepository) GetAllKelasMahasiswa(ctx context.Context, tx *gorm.DB, mahasiswaId uuid.UUID) ([]entities.KelasMahasiswa, error) {
	if tx == nil {
		return nil, nil
	}

	var entities []entities.KelasMahasiswa
	if err := tx.WithContext(ctx).
		Preload("Mahasiswa", helpers.SelectFields("detail_id, email, name")).
		Preload("Kelas", helpers.SelectFields("kelas_id, name")).
		Where("mahasiswa_id = ?", mahasiswaId).Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (k *kelasMahasiswaRepository) CheckMahasiswaAlreadyAssigned(ctx context.Context, tx *gorm.DB, mahasiswaId uuid.UUID) (bool, error) {
	if tx == nil {
		return false, nil
	}

	var count int64
	if err := tx.WithContext(ctx).
		Model(&entities.KelasMahasiswa{}).
		Where("mahasiswa_id = ? ", mahasiswaId).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
