package service

import (
	"context"
	"errors"
	"web-hosting/internal/database/entities"
	authDto "web-hosting/internal/modules/auth/dto"
	"web-hosting/internal/modules/auth/repository"
	userDto "web-hosting/internal/modules/user/dto"
	userRepo "web-hosting/internal/modules/user/repository"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/helpers"

	"gorm.io/gorm"
)

type AuthService interface {
	Login(ctx context.Context, req userDto.UserLoginRequest) (authDto.TokenResponse, error)
	Logout(ctx context.Context, userId string) error
	RefreshToken(ctx context.Context, req authDto.RefreshTokenRequest) (authDto.TokenResponse, error)
	ResetPassword(ctx context.Context, req authDto.ResetPasswordRequest) error
}
type authService struct {
	useRepo          userRepo.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtService       JwtService
	db               *gorm.DB
}

// func NewAuthService(useRepo userRepo.UserRepository, refreshTokenRepo repository.RefreshTokenRepository, jwtService JwtService) AuthService {
// 	return &authService{
// 		useRepo:          useRepo,
// 		refreshTokenRepo: refreshTokenRepo,
// 		jwtService:       jwtService,
// 	}
// }

func (s *authService) Login(ctx context.Context, req userDto.UserLoginRequest) (authDto.TokenResponse, error) {
	user, isExist, err := s.useRepo.CheckEmail(ctx, s.db, req.Email)
	if err != nil {
		return authDto.TokenResponse{}, constants.ErrInternalErr
	}
	if !isExist {
		return authDto.TokenResponse{}, userDto.ErrUserNotFound
	}
	isValid, err := helpers.CheckPasswordHash(req.Password, user.Password)
	if err != nil {
		return authDto.TokenResponse{}, constants.ErrInternalErr
	}
	if !isValid {
		return authDto.TokenResponse{}, authDto.ErrInvalidCredentials
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID.String(), user.Role.Name)
	if err != nil {
		return authDto.TokenResponse{}, constants.ErrInternalErr
	}

	refreshToken, exp := s.jwtService.GenerateRefreshToken()

	refreshTokenEntity := entities.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiredAt: exp,
	}

	_, err = s.refreshTokenRepo.Create(ctx, s.db, refreshTokenEntity)
	if err != nil {
		return authDto.TokenResponse{}, constants.ErrInternalErr
	}

	return authDto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		RoleName:     user.Role.Name,
	}, nil
}

func (s *authService) Logout(ctx context.Context, userId string) error {
	return s.refreshTokenRepo.DeleteByUserID(ctx, s.db, userId)
}

func (s *authService) RefreshToken(ctx context.Context, req authDto.RefreshTokenRequest) (authDto.TokenResponse, error) {
	refreshTokenEntity, err := s.refreshTokenRepo.FindByToken(ctx, s.db, req.RefreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return authDto.TokenResponse{}, authDto.ErrRefreshTokenNotFound
		}
		return authDto.TokenResponse{}, constants.ErrInternalErr
	}
	if err := s.refreshTokenRepo.DeleteByToken(ctx, s.db, req.RefreshToken); err != nil {
		return authDto.TokenResponse{}, constants.ErrInternalErr
	}

	accessToken, err := s.jwtService.GenerateAccessToken(refreshTokenEntity.UserID.String())
	if err != nil {
		return authDto.TokenResponse{}, constants.ErrInternalErr
	}

	refreshTokenNew, exp := s.jwtService.GenerateRefreshToken()

	refreshTokenEntityNew := entities.RefreshToken{
		UserID:    refreshTokenEntity.UserID,
		Token:     refreshTokenNew,
		ExpiredAt: exp,
	}

	_, err = s.refreshTokenRepo.Create(ctx, s.db, refreshTokenEntityNew)
	if err != nil {
		return authDto.TokenResponse{}, constants.ErrInternalErr
	}

	return authDto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenNew,
		RoleName:     refreshTokenEntity.User.Role.Name,
	}, nil
}

func (s *authService) ResetPassword(ctx context.Context, req authDto.ResetPasswordRequest) error {
	user, isExist, err := s.useRepo.CheckEmail(ctx, s.db, req.Email)
	if err != nil {
		return constants.ErrInternalErr
	}
	if !isExist {
		return userDto.ErrUserNotFound
	}

	_, err := s.useRepo.Update()
}
