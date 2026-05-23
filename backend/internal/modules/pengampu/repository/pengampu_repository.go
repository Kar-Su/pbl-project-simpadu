package repository

import (
	"context"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/package/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PengampuRepository interface {
		Create(ctx context.Context, tx *gorm.DB, entity entities.Pengampu) error
		UpdateByID(ctx context.Context, tx *gorm.DB, id uuid.UUID, entity entities.Pengampu) (entities.Pengampu, error)
		DeleteByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
		GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entities.Pengampu, error)
		GetByKelasID(ctx context.Context, tx *gorm.DB, kelasID uuid.UUID) ([]entities.Pengampu, error)
		GetByDosenID(ctx context.Context, tx *gorm.DB, dosenID uuid.UUID) ([]entities.Pengampu, error)
	}

	pengampuRepository struct {
		db *gorm.DB
	}
)

func NewPengampuRepository(db *gorm.DB) PengampuRepository {
	return &pengampuRepository{db: db}
}

func (r *pengampuRepository) Create(ctx context.Context, tx *gorm.DB, entity entities.Pengampu) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&entity).Error; err != nil {
		return err
	}

	return nil
}

func (r *pengampuRepository) UpdateByID(ctx context.Context, tx *gorm.DB, id uuid.UUID, entity entities.Pengampu) (entities.Pengampu, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Model(&entities.Pengampu{}).Where("id = ?", id).Omit("id").Updates(&entity).Error; err != nil {
		return entities.Pengampu{}, err
	}

	var updatedPengampu entities.Pengampu
	if err := tx.WithContext(ctx).Preload("MataKuliah").Preload("Dosen", helpers.SelectFields("detail_id, name, email")).Where("id = ?", id).First(&updatedPengampu).Error; err != nil {
		return entities.Pengampu{}, nil
	}

	return updatedPengampu, nil
}

func (r *pengampuRepository) DeleteByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	result := tx.WithContext(ctx).Preload("MataKuliah").Preload("Dosen", helpers.SelectFields("detail_id, name, email")).Where("id = ?", id).Delete(&entities.Pengampu{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *pengampuRepository) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entities.Pengampu, error) {
	if tx == nil {
		tx = r.db
	}

	var pengampu entities.Pengampu
	if err := tx.WithContext(ctx).Preload("MataKuliah").Preload("Dosen", helpers.SelectFields("detail_id, name, email")).Where("id = ?", id).First(&pengampu).Error; err != nil {
		return entities.Pengampu{}, err
	}

	return pengampu, nil
}

func (r *pengampuRepository) GetByKelasID(ctx context.Context, tx *gorm.DB, kelasID uuid.UUID) ([]entities.Pengampu, error) {
	if tx == nil {
		tx = r.db
	}

	var pengampu []entities.Pengampu
	if err := tx.WithContext(ctx).Preload("MataKuliah").Preload("Dosen", helpers.SelectFields("detail_id, name, email")).Where("kelas_id = ?", kelasID).Find(&pengampu).Error; err != nil {
		return nil, err
	}

	return pengampu, nil
}

func (r *pengampuRepository) GetByDosenID(ctx context.Context, tx *gorm.DB, dosenID uuid.UUID) ([]entities.Pengampu, error) {
	if tx == nil {
		tx = r.db
	}

	var pengampu []entities.Pengampu
	if err := tx.WithContext(ctx).Preload("MataKuliah").Preload("Dosen", helpers.SelectFields("detail_id, name, email")).Where("dosen_id = ?", dosenID).Find(&pengampu).Error; err != nil {
		return nil, err
	}

	return pengampu, nil
}
