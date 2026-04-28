package service

import (
	"context"
	"errors"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/modules/role/dto"
	"web-hosting/internal/modules/role/repository"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/helpers"

	"gorm.io/gorm"
)

type RoleService interface {
	Create(ctx context.Context, req dto.RoleCreateRequest) (entities.Role, error)
	Update(ctx context.Context, req dto.RoleUpdateRequest, roleId uint) (entities.Role, error)
	Delete(ctx context.Context, roleId uint) error
	GetRoleIdByRoleName(ctx context.Context, roleName string) (uint, error)
	GetRoleById(ctx context.Context, roleId uint) (entities.Role, error)
	GetAllRole(ctx context.Context) ([]entities.Role, error)
}

type roleService struct {
	roleRepo repository.RoleRepository
	db       *gorm.DB
}

func NewRoleService(roleRepo repository.RoleRepository) RoleService {
	return &roleService{roleRepo: roleRepo}
}

func (s *roleService) Create(ctx context.Context, req dto.RoleCreateRequest) (entities.Role, error) {
	normRoleName := helpers.NormalizeString(req.RoleName)
	isExist, err := s.roleRepo.CheckRoleExist(ctx, s.db, normRoleName)

	if err != nil {
		return entities.Role{}, constants.ErrInternalErr
	}

	if isExist {
		return entities.Role{}, dto.ErrRoleAlreadyExist
	}

	role, err := s.roleRepo.Create(ctx, s.db, normRoleName)
	if err != nil {
		return entities.Role{}, constants.ErrInternalErr
	}

	return role, nil
}

func (s *roleService) Update(ctx context.Context, req dto.RoleUpdateRequest, roleId uint) (entities.Role, error) {
	role, err := s.roleRepo.GetRoleById(ctx, s.db, roleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Role{}, dto.ErrRoleNotFound
		}
		return entities.Role{}, constants.ErrInternalErr
	}

	normRoleName := helpers.NormalizeString(req.RoleName)
	role.Name = normRoleName
	roleUpdated, err := s.roleRepo.Update(ctx, s.db, roleId, role)
	if err != nil {
		return entities.Role{}, constants.ErrInternalErr
	}

	return roleUpdated, nil
}

func (s *roleService) Delete(ctx context.Context, roleId uint) error {
	_, err := s.roleRepo.GetRoleById(ctx, s.db, roleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrRoleNotFound
		}
		return constants.ErrInternalErr
	}

	err = s.roleRepo.Delete(ctx, s.db, roleId)
	if err != nil {
		return constants.ErrInternalErr
	}

	return nil
}

func (s *roleService) GetRoleIdByRoleName(ctx context.Context, roleName string) (uint, error) {
	roleName = helpers.NormalizeString(roleName)
	roleId, err := s.roleRepo.GetRoleIdByRoleName(ctx, s.db, roleName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, dto.ErrRoleNotFound
		}
		return 0, constants.ErrInternalErr
	}

	return roleId, nil
}

func (s *roleService) GetRoleById(ctx context.Context, roleId uint) (entities.Role, error) {
	role, err := s.roleRepo.GetRoleById(ctx, s.db, roleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Role{}, dto.ErrRoleNotFound
		}
		return entities.Role{}, constants.ErrInternalErr
	}

	return role, nil
}

func (s *roleService) GetAllRole(ctx context.Context) ([]entities.Role, error) {
	roles, err := s.roleRepo.GetAllRole(ctx, s.db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrRoleNotFound
		}
		return nil, constants.ErrInternalErr
	}

	return roles, nil
}
