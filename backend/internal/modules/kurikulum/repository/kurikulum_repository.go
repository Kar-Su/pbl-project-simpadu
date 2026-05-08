package repository

import (
	"context"
	"errors"
	"web-hosting/internal/database/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	KurikulumRepository interface {
		Create(ctx context.Context, tx *gorm.DB, entity entities.Kurikulum) error
		UpdateByID(ctx context.Context, tx *gorm.DB, entity entities.Kurikulum, id uuid.UUID) error
		UpdateByKode(ctx context.Context, tx *gorm.DB, entity entities.Kurikulum, kode string) error
		DeleteByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
		DeleteByKode(ctx context.Context, tx *gorm.DB, kode string) error
		GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entities.Kurikulum, error)
		GetByKode(ctx context.Context, tx *gorm.DB, kode string) (entities.Kurikulum, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]entities.Kurikulum, error)
		CheckKurikulumExistsByKode(ctx context.Context, tx *gorm.DB, kode string) (bool, error)
	}

	kurikulumRepo struct {
		db *gorm.DB
	}
)

func NewKurikulumRepository(db *gorm.DB) KurikulumRepository {
	return &kurikulumRepo{db: db}
}

func (r *kurikulumRepo) Create(ctx context.Context, tx *gorm.DB, entity entities.Kurikulum) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&entity).Error; err != nil {
		return err
	}
	return nil
}

func (r *kurikulumRepo) UpdateByID(ctx context.Context, tx *gorm.DB, entity entities.Kurikulum, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Where("id = ?", id).Updates(&entity).Error; err != nil {
		return err
	}
	return nil
}

func (r *kurikulumRepo) UpdateByKode(ctx context.Context, tx *gorm.DB, entity entities.Kurikulum, kode string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Where("kode = ?", kode).Updates(&entity).Error; err != nil {
		return err
	}
	return nil
}

func (r *kurikulumRepo) DeleteByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	result := tx.WithContext(ctx).Where("id = ?", id).Delete(entities.Kurikulum{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *kurikulumRepo) DeleteByKode(ctx context.Context, tx *gorm.DB, kode string) error {
	if tx == nil {
		tx = r.db
	}

	result := tx.WithContext(ctx).Where("kode = ?", kode).Delete(entities.Kurikulum{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *kurikulumRepo) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entities.Kurikulum, error) {
	if tx == nil {
		tx = r.db
	}

	var entity entities.Kurikulum
	if err := tx.WithContext(ctx).Preload("KurikulumMK.MataKuliah").Preload("Prodi.Jurusan").Where("id = ?", id).First(&entity).Error; err != nil {
		return entities.Kurikulum{}, err
	}
	return entity, nil
}

func (r *kurikulumRepo) GetByKode(ctx context.Context, tx *gorm.DB, kode string) (entities.Kurikulum, error) {
	if tx == nil {
		tx = r.db
	}

	var entity entities.Kurikulum
	if err := tx.WithContext(ctx).Preload("KurikulumMK.MataKuliah").Preload("Prodi.Jurusan").Where("kode = ?", kode).First(&entity).Error; err != nil {
		return entities.Kurikulum{}, err
	}
	return entity, nil
}

func (r *kurikulumRepo) GetAll(ctx context.Context, tx *gorm.DB) ([]entities.Kurikulum, error) {
	if tx == nil {
		tx = r.db
	}

	var entities []entities.Kurikulum
	if err := tx.WithContext(ctx).Preload("KurikulumMK.MataKuliah").Preload("Prodi.Jurusan").Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *kurikulumRepo) CheckKurikulumExistsByKode(ctx context.Context, tx *gorm.DB, kode string) (bool, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Select("kode").Where("kode = ?", kode).First(&entities.Kurikulum{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
