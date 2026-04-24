package repository

import (
	"context"
	"errors"
	"web-hosting/internal/database/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Register(ctx context.Context, tx *gorm.DB, user entities.User) (entities.User, error)
	Update(ctx context.Context, tx *gorm.DB, userid uuid.UUID, user entities.User) (entities.User, error)
	UpdateByRoleAndDetailID(ctx context.Context, tx *gorm.DB, roleKode string, detailId uint, user entities.User) (entities.User, error)
	Delete(ctx context.Context, tx *gorm.DB, userId uuid.UUID) error
	DeleteByRoleAndDetailID(ctx context.Context, tx *gorm.DB, roleKode string, detailId uint) error
	GetUserByID(ctx context.Context, tx *gorm.DB, userId uuid.UUID) (entities.User, error)
	GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (entities.User, error)
	GetUserByRole(ctx context.Context, tx *gorm.DB, roleKode string) ([]entities.User, error)
	GetUserByRoleAndDetailID(ctx context.Context, tx *gorm.DB, roleKode string, detailId uint) (entities.User, error)
	CheckEmail(ctx context.Context, tx *gorm.DB, email string) (entities.User, bool, error)
	CheckRoleWithDetailID(ctx context.Context, tx *gorm.DB, roleKode string, detailId uint) (entities.User, bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Register(ctx context.Context, tx *gorm.DB, user entities.User) (entities.User, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, tx *gorm.DB, userId uuid.UUID, user entities.User) (entities.User, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Model(&entities.User{}).Where("id = ?", user.ID).Updates(&user).Error; err != nil {
		return entities.User{}, err
	}

	var updatedUser entities.User
	if err := tx.WithContext(ctx).First(&updatedUser, "id = ?", userId).Error; err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (r *userRepository) UpdateByRoleAndDetailID(ctx context.Context, tx *gorm.DB, roleKode string, detailId uint, user entities.User) (entities.User, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Model(&entities.User{}).Where("role_kode = ? AND detail_id = ?", roleKode, detailId).Select("*").Updates(&user).Error; err != nil {
		return entities.User{}, err
	}

	var updatedUser entities.User
	if err := tx.WithContext(ctx).First(&updatedUser, "role_kode = ? AND detail_id = ?", roleKode, detailId).Error; err != nil {
		return entities.User{}, err
	}

	return updatedUser, nil
}

func (r *userRepository) Delete(ctx context.Context, tx *gorm.DB, userId uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Delete(&entities.User{}, "id = ?", userId).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) DeleteByRoleAndDetailID(ctx context.Context, tx *gorm.DB, roleKode string, detailId uint) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Delete(&entities.User{}, "role_kode = ? AND detail_id = ?", roleKode, detailId).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetUserByID(ctx context.Context, tx *gorm.DB, userId uuid.UUID) (entities.User, error) {
	if tx == nil {
		tx = r.db
	}
	var user entities.User
	if err := tx.WithContext(ctx).First(&user, "id = ?", userId).Error; err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (entities.User, error) {
	if tx == nil {
		tx = r.db
	}
	var user entities.User
	if err := tx.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func (r *userRepository) GetUserByRole(ctx context.Context, tx *gorm.DB, roleKode string) ([]entities.User, error) {
	if tx == nil {
		tx = r.db
	}
	var users []entities.User
	if err := tx.WithContext(ctx).Find(&users, "role_kode = ?", roleKode).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetUserByRoleAndDetailID(ctx context.Context, tx *gorm.DB, roleKode string, detailId uint) (entities.User, error) {
	if tx == nil {
		tx = r.db
	}
	var users entities.User
	if err := tx.WithContext(ctx).Find(&users, "role_kode = ? AND detail_id = ?", roleKode, detailId).Error; err != nil {
		return entities.User{}, err
	}
	return users, nil
}

func (r *userRepository) CheckEmail(ctx context.Context, tx *gorm.DB, email string) (entities.User, bool, error) {
	if tx == nil {
		tx = r.db
	}
	var user entities.User
	if err := tx.WithContext(ctx).Where("email = ?", email).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, false, nil
		}
		return entities.User{}, false, err
	}
	return user, true, nil
}

func (r *userRepository) CheckRoleWithDetailID(ctx context.Context, tx *gorm.DB, roleKode string, detailId uint) (entities.User, bool, error) {
	if tx == nil {
		tx = r.db
	}
	var user entities.User
	if err := tx.WithContext(ctx).First(&user, "role_kode = ? AND detail_id = ?", roleKode, detailId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, false, nil
		}

		return entities.User{}, false, err
	}
	return user, true, nil
}
