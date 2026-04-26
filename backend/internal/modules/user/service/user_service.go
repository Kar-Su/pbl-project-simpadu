package service

import (
	"context"
	"errors"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/modules/user/dto"
	"web-hosting/internal/modules/user/repository"
	"web-hosting/internal/package/constants"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	CreateAdmin(ctx context.Context, req dto.UserAdminCreateRequest) (dto.UserResponse, error)
	CreateNonAdmin(ctx context.Context, req dto.UserNonAdminCreateRequest) (dto.UserResponse, error)
	UpdateAdmin(ctx context.Context, req dto.UserAdminUpdateRequest, userId uuid.UUID) (dto.UserResponse, error)
	UpdateNonAdmin(ctx context.Context, req dto.UserNonAdminUpdateRequest, roleId uint, detailId uint) (dto.UserResponse, error)
	DeleteAdmin(ctx context.Context, userId uuid.UUID) error
	DeleteNonAdmin(ctx context.Context, roleId uint, detailId uint) error
	GetUserByID(ctx context.Context, userId uuid.UUID) (dto.UserResponse, error)
	GetUserByRoleAndDetailID(ctx context.Context, roleId uint, detailId uint) (dto.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, error)
	GetUserByRole(ctx context.Context, roleId uint) ([]dto.UserResponse, error)
}

type userService struct {
	userRepository repository.UserRepository
	db             *gorm.DB
}

func NewUserService(userRepository repository.UserRepository, db *gorm.DB) UserService {
	return &userService{
		userRepository: userRepository,
		db:             db,
	}
}

func (s *userService) CreateAdmin(ctx context.Context, req dto.UserAdminCreateRequest) (dto.UserResponse, error) {
	_, isExist, err := s.userRepository.CheckEmail(ctx, s.db, req.Email)
	if err != nil {
		return dto.UserResponse{}, constants.ErrInternalErr
	}
	if isExist {
		return dto.UserResponse{}, dto.ErrEmailAlreadyExists
	}

	userEntity := entities.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		RoleID:   dto.RoleNameToRoleID(req.RoleName),
	}
	if req.DetailId != nil {
		userEntity.DetailID = req.DetailId
	}
	if req.Image != nil {
		fileName := req.Image.Filename
		userEntity.ImageUrl = &fileName
	}
	userCreated, err := s.userRepository.Register(ctx, s.db, userEntity)
	if err != nil {
		return dto.UserResponse{}, dto.ErrCreateUser
	}

	return dto.ToUserResponse(userCreated), nil
}

func (s *userService) CreateNonAdmin(ctx context.Context, req dto.UserNonAdminCreateRequest) (dto.UserResponse, error) {
	_, isExist, err := s.userRepository.CheckEmail(ctx, s.db, req.Email)
	if err != nil {
		return dto.UserResponse{}, constants.ErrInternalErr
	}
	if isExist {
		return dto.UserResponse{}, dto.ErrEmailAlreadyExists
	}

	userEntity := entities.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		RoleID:   dto.RoleNameToRoleID(req.RoleName),
	}
	if req.DetailId != nil {
		userEntity.DetailID = req.DetailId
	}
	if req.Image != nil {
		fileName := req.Image.Filename
		userEntity.ImageUrl = &fileName
	}
	userCreated, err := s.userRepository.Register(ctx, s.db, userEntity)
	if err != nil {
		return dto.UserResponse{}, dto.ErrCreateUser
	}

	return dto.ToUserResponse(userCreated), nil
}

func (s *userService) UpdateAdmin(ctx context.Context, req dto.UserAdminUpdateRequest, userId uuid.UUID) (dto.UserResponse, error) {
	user, err := s.userRepository.GetUserByID(ctx, s.db, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserResponse{}, dto.ErrUserNotFound
		}
		return dto.UserResponse{}, constants.ErrInternalErr
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		user.Password = req.Password
	}
	if roleName := req.RoleName; roleName != "" {
		user.RoleID = dto.RoleNameToRoleID(roleName)
	}
	if req.DetailId != nil {
		user.DetailID = req.DetailId
	}
	if req.Image != nil {
		fileName := req.Image.Filename
		user.ImageUrl = &fileName
	}
	updatedUser, err := s.userRepository.Update(ctx, s.db, userId, user)
	if err != nil {
		return dto.UserResponse{}, dto.ErrUpdateUser
	}

	return dto.ToUserResponse(updatedUser), nil
}

func (s *userService) UpdateNonAdmin(ctx context.Context, req dto.UserNonAdminUpdateRequest, roleId uint, detailId uint) (dto.UserResponse, error) {
	user, isExist, err := s.userRepository.CheckRoleWithDetailID(ctx, s.db, roleId, detailId)
	if err != nil {
		return dto.UserResponse{}, constants.ErrInternalErr
	}
	if !isExist {
		return dto.UserResponse{}, dto.ErrUserNotFound
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		user.Password = req.Password
	}
	if req.Image != nil {
		if req.Image.Filename == "null" {
			user.ImageUrl = nil
		} else {
			fileName := req.Image.Filename
			user.ImageUrl = &fileName
		}
	}
	updatedUser, err := s.userRepository.UpdateByRoleAndDetailID(ctx, s.db, roleId, detailId, user)
	if err != nil {
		return dto.UserResponse{}, dto.ErrUpdateUser
	}

	return dto.ToUserResponse(updatedUser), nil
}

func (s *userService) DeleteAdmin(ctx context.Context, userId uuid.UUID) error {
	_, err := s.userRepository.GetUserByID(ctx, s.db, userId)
	if err != nil {
		if errors.Is(err, dto.ErrUserNotFound) {
			return dto.ErrUserNotFound
		}
		return constants.ErrInternalErr
	}

	if err := s.userRepository.Delete(ctx, s.db, userId); err != nil {
		return dto.ErrDeleteUser
	}
	return nil
}

func (s *userService) DeleteNonAdmin(ctx context.Context, roleId uint, detailId uint) error {
	user, isExist, err := s.userRepository.CheckRoleWithDetailID(ctx, s.db, roleId, detailId)
	if err != nil {
		return constants.ErrInternalErr
	}
	if !isExist {
		return dto.ErrUserNotFound
	}
	if err := s.userRepository.Delete(ctx, s.db, user.ID); err != nil {
		return dto.ErrDeleteUser
	}
	return nil
}

func (s *userService) GetUserByID(ctx context.Context, userId uuid.UUID) (dto.UserResponse, error) {
	user, err := s.userRepository.GetUserByID(ctx, s.db, userId)
	if err != nil {
		if errors.Is(err, dto.ErrUserNotFound) {
			return dto.UserResponse{}, dto.ErrUserNotFound
		}
		return dto.UserResponse{}, constants.ErrInternalErr
	}
	return dto.ToUserResponse(user), nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, s.db, email)
	if err != nil {
		if errors.Is(err, dto.ErrUserNotFound) {
			return dto.UserResponse{}, dto.ErrUserNotFound
		}
		return dto.UserResponse{}, constants.ErrInternalErr
	}
	return dto.ToUserResponse(user), nil
}

func (s *userService) GetUserByRole(ctx context.Context, roleId uint) ([]dto.UserResponse, error) {
	user, err := s.userRepository.GetUserByRole(ctx, s.db, roleId)
	if err != nil {
		return nil, dto.ErrUserNotFound
	}

	responses := make([]dto.UserResponse, 0, len(user))
	for _, u := range user {
		responses = append(responses, dto.ToUserResponse(u))
	}

	return responses, nil
}

func (s *userService) GetUserByRoleAndDetailID(ctx context.Context, roleId uint, detailId uint) (dto.UserResponse, error) {
	user, err := s.userRepository.GetUserByRoleAndDetailID(ctx, s.db, roleId, detailId)
	if err != nil {
		if errors.Is(err, dto.ErrUserNotFound) {
			return dto.UserResponse{}, dto.ErrUserNotFound
		}
		return dto.UserResponse{}, constants.ErrInternalErr
	}
	return dto.ToUserResponse(user), nil
}
