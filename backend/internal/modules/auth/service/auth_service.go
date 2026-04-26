package service

import (
	"context"
	authDto "web-hosting/internal/modules/auth/dto"
	"web-hosting/internal/modules/auth/repository"
	userDto "web-hosting/internal/modules/user/dto"
	userRepo "web-hosting/internal/modules/user/repository"

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
	_, isExist, err := s.useRepo.CheckEmail(ctx, s.db, req.Email)
	if err != nil {

	}
}
