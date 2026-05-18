package repository

import (
	"context"
	"errors"
	"web-hosting/internal/database/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	MkRepository interface {
		Create(ctx context.Context, tx *gorm.DB, mk entities.MataKuliah) error
		UpdateByKode(ctx context.Context, tx *gorm.DB, kodeMk string, mk entities.MataKuliah) (entities.MataKuliah, error)
		UpdateById(ctx context.Context, tx *gorm.DB, id uuid.UUID, mk entities.MataKuliah) (entities.MataKuliah, error)
		DeleteByKode(ctx context.Context, tx *gorm.DB, kodeMk string) error
		DeleteById(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
		GetByKode(ctx context.Context, tx *gorm.DB, kodeMk string) (entities.MataKuliah, error)
		GetById(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entities.MataKuliah, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]entities.MataKuliah, error)
		GetAllPaginated(ctx context.Context, tx *gorm.DB, offset, limit int) ([]entities.MataKuliah, int64, error)
		CheckMkExists(ctx context.Context, tx *gorm.DB, kodeMk string) (bool, error)
	}

	mkRepository struct {
		db *gorm.DB
	}
)

func NewMkRepository(db *gorm.DB) MkRepository {
	return &mkRepository{db: db}
}

func (r *mkRepository) Create(ctx context.Context, tx *gorm.DB, mk entities.MataKuliah) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&mk).Error; err != nil {
		return err
	}

	return nil
}

func (r *mkRepository) UpdateByKode(ctx context.Context, tx *gorm.DB, kodeMk string, mk entities.MataKuliah) (entities.MataKuliah, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Where("kode = ?", kodeMk).Updates(&mk).Error; err != nil {
		return entities.MataKuliah{}, err
	}

	updated := entities.MataKuliah{}
	if err := tx.WithContext(ctx).Where("kode = ?", kodeMk).First(&updated).Error; err != nil {
		return entities.MataKuliah{}, err
	}

	return updated, nil
}

func (r *mkRepository) UpdateById(ctx context.Context, tx *gorm.DB, id uuid.UUID, mk entities.MataKuliah) (entities.MataKuliah, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Where("id = ?", id).Updates(&mk).Error; err != nil {
		return entities.MataKuliah{}, err
	}

	updated := entities.MataKuliah{}
	if err := tx.WithContext(ctx).Where("id = ?", id).First(&updated).Error; err != nil {
		return entities.MataKuliah{}, err
	}

	return updated, nil
}

func (r *mkRepository) DeleteByKode(ctx context.Context, tx *gorm.DB, kodeMk string) error {
	if tx == nil {
		tx = r.db
	}

	res := tx.WithContext(ctx).Where("kode = ?", kodeMk).Delete(&entities.MataKuliah{})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *mkRepository) DeleteById(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	res := tx.WithContext(ctx).Where("id = ?", id).Delete(&entities.MataKuliah{})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *mkRepository) GetByKode(ctx context.Context, tx *gorm.DB, kodeMk string) (entities.MataKuliah, error) {
	if tx == nil {
		tx = r.db
	}

	mk := entities.MataKuliah{}
	if err := tx.WithContext(ctx).Where("kode = ?", kodeMk).First(&mk).Error; err != nil {
		return entities.MataKuliah{}, err
	}

	return mk, nil
}

func (r *mkRepository) GetById(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entities.MataKuliah, error) {
	if tx == nil {
		tx = r.db
	}

	mk := entities.MataKuliah{}
	if err := tx.WithContext(ctx).Where("id = ?", id).First(&mk).Error; err != nil {
		return entities.MataKuliah{}, err
	}

	return mk, nil
}

func (r *mkRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entities.MataKuliah, error) {
	if tx == nil {
		tx = r.db
	}

	var mks []entities.MataKuliah
	if err := tx.WithContext(ctx).Find(&mks).Error; err != nil {
		return nil, err
	}

	return mks, nil
}

func (r *mkRepository) GetAllPaginated(ctx context.Context, tx *gorm.DB, offset, limit int) ([]entities.MataKuliah, int64, error) {
	if tx == nil {
		tx = r.db
	}
	var mks []entities.MataKuliah
	var total int64
	if err := tx.WithContext(ctx).Model(&entities.MataKuliah{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.WithContext(ctx).Offset(offset).Limit(limit).Find(&mks).Error; err != nil {
		return nil, 0, err
	}
	return mks, total, nil
}

func (r *mkRepository) CheckMkExists(ctx context.Context, tx *gorm.DB, kodeMk string) (bool, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Model(&entities.MataKuliah{}).Select("kode").Where("kode = ?", kodeMk).First(&entities.MataKuliah{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
