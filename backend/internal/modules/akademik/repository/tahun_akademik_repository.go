package repository

import (
	"context"
	"web-hosting/internal/database/entities"

	"gorm.io/gorm"
)

type (
	TahunAkademikRepository interface {
		Create(ctx context.Context, tx *gorm.DB, entity entities.TahunAkademik) error
		Update(ctx context.Context, tx *gorm.DB, id uint, entity entities.TahunAkademik) error
		Delete(ctx context.Context, tx *gorm.DB, id uint) error
		GetByID(ctx context.Context, tx *gorm.DB, id uint) (entities.TahunAkademik, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]entities.TahunAkademik, error)
		GetByStatus(ctx context.Context, tx *gorm.DB, status string) ([]entities.TahunAkademik, error)
		CheckByID(ctx context.Context, tx *gorm.DB, id uint) (bool, error)
	}
	tahunAkademikRepository struct {
		db *gorm.DB
	}
)

func NewTahunAkademikRepository(db *gorm.DB) TahunAkademikRepository {
	return &tahunAkademikRepository{db: db}
}

func (r *tahunAkademikRepository) Create(ctx context.Context, tx *gorm.DB, entity entities.TahunAkademik) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&entity).Error; err != nil {
		return err
	}
	return nil
}

func (r *tahunAkademikRepository) Update(ctx context.Context, tx *gorm.DB, id uint, entity entities.TahunAkademik) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Where("id = ?", id).Updates(&entity).Error; err != nil {
		return err
	}

	return nil
}

func (r *tahunAkademikRepository) Delete(ctx context.Context, tx *gorm.DB, id uint) error {
	if tx == nil {
		tx = r.db
	}

	result := tx.WithContext(ctx).Where("id = ?", id).Delete(&entities.TahunAkademik{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *tahunAkademikRepository) GetByID(ctx context.Context, tx *gorm.DB, id uint) (entities.TahunAkademik, error) {
	if tx == nil {
		tx = r.db
	}

	var entity entities.TahunAkademik
	if err := tx.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
		return entities.TahunAkademik{}, err
	}

	return entity, nil
}

func (r *tahunAkademikRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entities.TahunAkademik, error) {
	if tx == nil {
		tx = r.db
	}

	var entities []entities.TahunAkademik
	if err := tx.WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *tahunAkademikRepository) GetByStatus(ctx context.Context, tx *gorm.DB, status string) ([]entities.TahunAkademik, error) {
	if tx == nil {
		tx = r.db
	}

	var entities []entities.TahunAkademik
	if err := tx.WithContext(ctx).Where("status = ?", status).Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *tahunAkademikRepository) CheckByID(ctx context.Context, tx *gorm.DB, id uint) (bool, error) {
	if tx == nil {
		tx = r.db
	}

	var entity entities.TahunAkademik
	if err := tx.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
