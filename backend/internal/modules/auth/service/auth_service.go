package service

import (
	"context"
	"errors"
	"log"
	"time"
	"web-hosting/internal/database/entities"
	authDto "web-hosting/internal/modules/auth/dto"
	"web-hosting/internal/modules/auth/repository"
	kelasRepo "web-hosting/internal/modules/kelas/repository"
	userDto "web-hosting/internal/modules/user/dto"
	userRepo "web-hosting/internal/modules/user/repository"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService interface {
	FindRefreshToken(ctx context.Context, token string) (authDto.RefreshTokenResponse, error)
	Login(ctx context.Context, req userDto.UserLoginRequest) (authDto.TokenResponse, error)
	Logout(ctx context.Context, userId string) error
	RefreshToken(ctx context.Context, req authDto.RefreshTokenRequest) (authDto.TokenResponse, error)
	ResetPassword(ctx context.Context, req authDto.ResetPasswordRequest) error
	CleanupExpiredTokens(ctx context.Context) error
}

type authService struct {
	useRepo          userRepo.UserRepository
	kelasPivotRepo   kelasRepo.KelasMahasiswaRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtService       JwtService
	db               *gorm.DB
}

func NewAuthService(useRepo userRepo.UserRepository, refreshTokenRepo repository.RefreshTokenRepository, kelasPivotRepo kelasRepo.KelasMahasiswaRepository, jwtService JwtService, db *gorm.DB) AuthService {
	return &authService{
		useRepo:          useRepo,
		kelasPivotRepo:   kelasPivotRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
		db:               db,
	}
}

func (s *authService) FindRefreshToken(ctx context.Context, token string) (authDto.RefreshTokenResponse, error) {
	refreshToken, err := s.refreshTokenRepo.FindByToken(ctx, s.db, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return authDto.RefreshTokenResponse{}, authDto.ErrRefreshTokenNotFound
		}
		log.Printf("Internal Error: %v", err)
		return authDto.RefreshTokenResponse{}, constants.ErrInternalErr
	}

	return authDto.RefreshTokenResponse{
		RefreshToken: refreshToken.Token,
		Exp:          refreshToken.ExpiredAt.Unix(),
	}, nil
}

func (s *authService) Login(ctx context.Context, req userDto.UserLoginRequest) (authDto.TokenResponse, error) {
	user, isExist, err := s.useRepo.CheckEmail(ctx, s.db, req.Email)
	if err != nil {
		log.Printf("Internal Error(%v): %v", req.Email, err)
		return authDto.TokenResponse{}, constants.ErrInternalErr
	}
	if !isExist {
		return authDto.TokenResponse{}, userDto.ErrUserNotFound
	}
	isValid, err := helpers.CheckPasswordHash(req.Password, user.Password)
	if err != nil && !isValid {
		return authDto.TokenResponse{}, err
	}

	kelasId := uuid.Nil
	if user.Role.Name == constants.ROLE_MAHASISWA {
		kelasId, err = s.kelasPivotRepo.GetKelasIdByMahasiswa(ctx, s.db, user.DetailID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
			}
			log.Printf("Internal Error: %v", err)
			return authDto.TokenResponse{}, constants.ErrInternalErr
		}
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID.String(), user.Role.Name, user.Email, user.DetailID, kelasId)
	if err != nil {
		log.Printf("Internal Error: %v", err)
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
		log.Printf("Internal Error: %v", err)
		return authDto.TokenResponse{}, constants.ErrInternalErr
	}

	return authDto.TokenResponse{
		UserName:     user.Name,
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
		log.Printf("Internal Error: %v", err)
		return authDto.TokenResponse{}, constants.ErrInternalErr
	}

	if refreshTokenEntity.ExpiredAt.Before(time.Now()) {
		s.refreshTokenRepo.DeleteByToken(ctx, s.db, req.RefreshToken)
		return authDto.TokenResponse{}, authDto.ErrRefreshTokenExpired
	}

	if err := s.refreshTokenRepo.DeleteByToken(ctx, s.db, req.RefreshToken); err != nil {
		log.Printf("Internal Error: %v", err)
		return authDto.TokenResponse{}, constants.ErrInternalErr
	}

	kelasId, err := s.kelasPivotRepo.GetKelasIdByMahasiswa(ctx, s.db, refreshTokenEntity.User.DetailID)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return authDto.TokenResponse{}, constants.ErrInternalErr
	}

	accessToken, err := s.jwtService.GenerateAccessToken(refreshTokenEntity.UserID.String(), refreshTokenEntity.User.Role.Name, refreshTokenEntity.User.Email, refreshTokenEntity.User.DetailID, kelasId)
	if err != nil {
		log.Printf("Internal Error: %v", err)
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
		log.Printf("Internal Error: %v", err)
		return authDto.TokenResponse{}, constants.ErrInternalErr
	}

	return authDto.TokenResponse{
		UserName:     refreshTokenEntity.User.Name,
		AccessToken:  accessToken,
		RefreshToken: refreshTokenNew,
		RoleName:     refreshTokenEntity.User.Role.Name,
	}, nil
}

func (s *authService) ResetPassword(ctx context.Context, req authDto.ResetPasswordRequest) error {
	user, isExist, err := s.useRepo.CheckEmail(ctx, s.db, req.Email)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	if !isExist {
		return userDto.ErrUserNotFound
	}

	hashPass, err := helpers.HashPassword(req.NewPassword)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	user.Password = hashPass
	_, err = s.useRepo.Update(ctx, s.db, user.ID, user)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	return nil
}

func (s *authService) CleanupExpiredTokens(ctx context.Context) error {
	err := s.refreshTokenRepo.DeleteExpired(ctx, nil)
	if err != nil {
		log.Printf("Gagal menghapus token kadaluarsa: %v", err)
		return err
	}
	log.Println("Berhasil membersihkan token yang kadaluarsa.")
	return nil
}
