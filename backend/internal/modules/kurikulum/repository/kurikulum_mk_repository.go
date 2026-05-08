package repository

import (
	"context"
	"web-hosting/internal/database/entities"

	"gorm.io/gorm"
)

type (
	KurikulumMKRepo interface {
		Create(ctx context.Context, tx *gorm.DB, entity entities.KurikulumMK) error
		Update(ctx context.Context, tx *gorm.DB, kurikulumKode string, MkKode string, entity entities.KurikulumMK) error
		Delete(ctx context.Context, tx *gorm.DB, kurikulumKode string, MkKode string) error
		Get(ctx context.Context, tx *gorm.DB, kurikulumKode string, MkKode string) (entities.KurikulumMK, error)
	}

	kurikulumMKRepo struct {
		db *gorm.DB
	}
)

func NewKurikulumMKRepository(db *gorm.DB) KurikulumMKRepo {
	return &kurikulumMKRepo{db: db}
}

func (r *kurikulumMKRepo) Create(ctx context.Context, tx *gorm.DB, entity entities.KurikulumMK) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&entity).Error; err != nil {
		return err
	}
	return nil
}

func (r *kurikulumMKRepo) Update(ctx context.Context, tx *gorm.DB, kurikulumKode string, MkKode string, entity entities.KurikulumMK) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Where("kurikulum_kode = ? AND mk_kode = ?", kurikulumKode, MkKode).Updates(&entity).Error; err != nil {
		return err
	}
	return nil
}

func (r *kurikulumMKRepo) Delete(ctx context.Context, tx *gorm.DB, kurikulumKode string, MkKode string) error {
	if tx == nil {
		tx = r.db
	}

	result := tx.WithContext(ctx).Where("kurikulum_kode = ? AND mk_kode = ?", kurikulumKode, MkKode).Delete(&entities.KurikulumMK{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return nil
	}

	return nil
}

func (r *kurikulumMKRepo) Get(ctx context.Context, tx *gorm.DB, kurikulumKode string, MkKode string) (entities.KurikulumMK, error) {
	if tx == nil {
		tx = r.db
	}

	var entity entities.KurikulumMK
	if err := tx.WithContext(ctx).Preload("MataKuliah").Where("kurikulum_kode = ? AND mk_kode = ?", kurikulumKode, MkKode).First(&entity).Error; err != nil {
		return entities.KurikulumMK{}, err
	}

	return entity, nil
}
