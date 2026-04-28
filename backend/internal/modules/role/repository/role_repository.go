package repository

import (
	"context"
	"errors"
	"web-hosting/internal/database/entities"

	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(ctx context.Context, tx *gorm.DB, roleName string) (entities.Role, error)
	Update(ctx context.Context, tx *gorm.DB, roleId uint, role entities.Role) (entities.Role, error)
	Delete(ctx context.Context, tx *gorm.DB, roleId uint) error
	GetRoleById(ctx context.Context, tx *gorm.DB, roleId uint) (entities.Role, error)
	GetAllRole(ctx context.Context, tx *gorm.DB) ([]entities.Role, error)
	GetRoleIdByRoleName(ctx context.Context, tx *gorm.DB, roleName string) (uint, error)
	CheckRoleExist(ctx context.Context, tx *gorm.DB, roleName string) (bool, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, tx *gorm.DB, roleName string) (entities.Role, error) {
	if tx == nil {
		tx = r.db
	}
	role := entities.Role{Name: roleName}
	if err := tx.WithContext(ctx).Create(&role).Error; err != nil {
		return entities.Role{}, err
	}
	return role, nil
}

func (r *roleRepository) Update(ctx context.Context, tx *gorm.DB, roleId uint, role entities.Role) (entities.Role, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Where("id = ?", roleId).Updates(&role).Error; err != nil {
		return entities.Role{}, err
	}
	return role, nil
}

func (r *roleRepository) Delete(ctx context.Context, tx *gorm.DB, roleId uint) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Where("id = ?", roleId).Delete(&entities.Role{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *roleRepository) GetRoleById(ctx context.Context, tx *gorm.DB, roleId uint) (entities.Role, error) {
	if tx == nil {
		tx = r.db
	}
	var role entities.Role
	if err := tx.WithContext(ctx).First(&role, roleId).Error; err != nil {
		return entities.Role{}, err
	}
	return role, nil
}

func (r *roleRepository) GetAllRole(ctx context.Context, tx *gorm.DB) ([]entities.Role, error) {
	if tx == nil {
		tx = r.db
	}
	var roles []entities.Role
	if err := tx.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) GetRoleIdByRoleName(ctx context.Context, tx *gorm.DB, roleName string) (uint, error) {
	if tx == nil {
		tx = r.db
	}
	var role entities.Role
	if err := tx.WithContext(ctx).Where("name = ?", roleName).First(&role).Error; err != nil {
		return 0, err
	}
	return role.ID, nil
}

func (r *roleRepository) CheckRoleExist(ctx context.Context, tx *gorm.DB, roleName string) (bool, error) {
	if tx == nil {
		tx = r.db
	}
	var role entities.Role
	if err := tx.WithContext(ctx).Where("name = ?", roleName).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
