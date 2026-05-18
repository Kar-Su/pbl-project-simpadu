package repository

import (
	"context"
	"web-hosting/internal/database/entities"

	"gorm.io/gorm"
)

type JurusanRepository interface {
	Create(ctx context.Context, tx *gorm.DB, jurusanName string) error
	Update(ctx context.Context, tx *gorm.DB, jurusanName string, jurusan entities.Jurusan) (entities.Jurusan, error)
	Delete(ctx context.Context, tx *gorm.DB, jurusanName string) error
	GetByName(ctx context.Context, tx *gorm.DB, jurusanName string) (entities.Jurusan, error)
	GetAll(ctx context.Context, tx *gorm.DB) ([]entities.Jurusan, error)
	GetAllPaginated(ctx context.Context, tx *gorm.DB, offset, limit int) ([]entities.Jurusan, int64, error)
	GetById(ctx context.Context, tx *gorm.DB, id uint) (entities.Jurusan, error)
}

type jurusanRepository struct {
	db *gorm.DB
}

func NewJurusanRepository(db *gorm.DB) JurusanRepository {
	return &jurusanRepository{db: db}
}

func (r *jurusanRepository) Create(ctx context.Context, tx *gorm.DB, jurusanName string) error {
	if tx == nil {
		tx = r.db
	}

	jurusan := entities.Jurusan{Name: jurusanName}
	if err := tx.WithContext(ctx).Create(&jurusan).Error; err != nil {
		return err
	}
	return nil
}

func (r *jurusanRepository) Update(ctx context.Context, tx *gorm.DB, jurusanName string, jurusan entities.Jurusan) (entities.Jurusan, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Where("name = ?", jurusanName).Updates(&jurusan).Error; err != nil {
		return entities.Jurusan{}, err
	}

	// Gunakan nama BARU (jurusan.Name) untuk fetch hasil update
	updatedJurusan, err := r.GetByName(ctx, tx, jurusan.Name)
	if err != nil {
		return entities.Jurusan{}, err
	}
	return updatedJurusan, nil
}

func (r *jurusanRepository) Delete(ctx context.Context, tx *gorm.DB, jurusanName string) error {
	if tx == nil {
		tx = r.db
	}
	result := tx.WithContext(ctx).Where("name = ?", jurusanName).Delete(&entities.Jurusan{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *jurusanRepository) GetByName(ctx context.Context, tx *gorm.DB, jurusanName string) (entities.Jurusan, error) {
	if tx == nil {
		tx = r.db
	}
	var jurusan entities.Jurusan
	if err := tx.WithContext(ctx).Where("name = ?", jurusanName).First(&jurusan).Error; err != nil {
		return entities.Jurusan{}, err
	}
	return jurusan, nil
}

func (r *jurusanRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entities.Jurusan, error) {
	if tx == nil {
		tx = r.db
	}
	var jurusan []entities.Jurusan
	if err := tx.WithContext(ctx).Find(&jurusan).Error; err != nil {
		return nil, err
	}
	return jurusan, nil
}

func (r *jurusanRepository) GetAllPaginated(ctx context.Context, tx *gorm.DB, offset, limit int) ([]entities.Jurusan, int64, error) {
	if tx == nil {
		tx = r.db
	}
	var jurusan []entities.Jurusan
	var total int64
	if err := tx.WithContext(ctx).Model(&entities.Jurusan{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.WithContext(ctx).Offset(offset).Limit(limit).Find(&jurusan).Error; err != nil {
		return nil, 0, err
	}
	return jurusan, total, nil
}

func (r *jurusanRepository) GetById(ctx context.Context, tx *gorm.DB, id uint) (entities.Jurusan, error) {
	if tx == nil {
		tx = r.db
	}
	var jurusan entities.Jurusan
	if err := tx.WithContext(ctx).Where("id = ?", id).First(&jurusan).Error; err != nil {
		return entities.Jurusan{}, err
	}
	return jurusan, nil
}
