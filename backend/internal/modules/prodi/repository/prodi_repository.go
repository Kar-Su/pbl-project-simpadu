package repository

import (
	"context"
	"web-hosting/internal/database/entities"

	"gorm.io/gorm"
)

type (
	ProdiRepository interface {
		Create(ctx context.Context, tx *gorm.DB, prodi entities.Prodi) error
		Update(ctx context.Context, tx *gorm.DB, prodiName string, prodi entities.Prodi) (entities.Prodi, error)
		Delete(ctx context.Context, tx *gorm.DB, prodiName string) error
		GetByName(ctx context.Context, tx *gorm.DB, prodiName string) (entities.Prodi, error)
		GetByID(ctx context.Context, tx *gorm.DB, prodiId uint) (entities.Prodi, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]entities.Prodi, error)
		GetByJurusanID(ctx context.Context, tx *gorm.DB, jurusanID uint) ([]entities.Prodi, error)
	}

	prodiRepository struct {
		db *gorm.DB
	}
)

func NewProdiRepository(db *gorm.DB) ProdiRepository {
	return &prodiRepository{db: db}
}

func (r *prodiRepository) Create(ctx context.Context, tx *gorm.DB, prodi entities.Prodi) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&prodi).Error; err != nil {
		return err
	}

	return nil
}

func (r *prodiRepository) Update(ctx context.Context, tx *gorm.DB, prodiName string, prodi entities.Prodi) (entities.Prodi, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Where("name = ?", prodiName).Updates(&prodi).Error; err != nil {
		return entities.Prodi{}, err
	}

	updatedProdi, err := r.GetByName(ctx, tx, prodiName)
	if err != nil {
		return entities.Prodi{}, err
	}

	return updatedProdi, nil
}

func (r *prodiRepository) Delete(ctx context.Context, tx *gorm.DB, prodiName string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Where("name = ?", prodiName).Delete(&entities.Prodi{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *prodiRepository) GetByName(ctx context.Context, tx *gorm.DB, prodiName string) (entities.Prodi, error) {
	if tx == nil {
		tx = r.db
	}

	var prodi entities.Prodi
	if err := tx.WithContext(ctx).Preload("Jurusan").Where("name = ?", prodiName).First(&prodi).Error; err != nil {
		return entities.Prodi{}, err
	}

	return prodi, nil
}

func (r *prodiRepository) GetByID(ctx context.Context, tx *gorm.DB, prodiID uint) (entities.Prodi, error) {
	if tx == nil {
		tx = r.db
	}

	var prodi entities.Prodi
	if err := tx.WithContext(ctx).Preload("Jurusan").Where("id = ?", prodiID).First(&prodi).Error; err != nil {
		return entities.Prodi{}, err
	}

	return prodi, nil
}

func (r *prodiRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entities.Prodi, error) {
	if tx == nil {
		tx = r.db
	}

	var prodis []entities.Prodi
	if err := tx.WithContext(ctx).Preload("Jurusan").Find(&prodis).Error; err != nil {
		return nil, err
	}

	return prodis, nil
}

func (r *prodiRepository) GetByJurusanID(ctx context.Context, tx *gorm.DB, jurusanID uint) ([]entities.Prodi, error) {
	if tx == nil {
		tx = r.db
	}

	var prodis []entities.Prodi
	if err := tx.WithContext(ctx).Preload("Jurusan").Where("jurusan_id = ?", jurusanID).Find(&prodis).Error; err != nil {
		return nil, err
	}

	return prodis, nil
}
