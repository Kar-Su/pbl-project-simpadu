package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"
	"web-hosting/internal/modules/auth/dto"
	"web-hosting/internal/package/constants"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	GenerateAccessToken(userId string, roleName string, userEmail string) (string, error)
	GenerateRefreshToken() (string, time.Time)
	ValidateToken(token string) (*jwt.Token, error)
	GetUserIDByToken(token string) (string, error)
	GetRoleNameByToken(token string) (string, error)
	GetUserEmailByToken(tokenString string) (string, error)
}

type jwtCustomClaim struct {
	UserID    string `json:"user_id"`
	RoleName  string `json:"role_name"`
	UserEmail string `json:"user_email"`
	jwt.RegisteredClaims
}

func getSecretKey() string {
	secretKey := constants.JWT_SECRET_KEY
	if secretKey == "" {
		panic(dto.MESSAGE_FAILED_RETRIEVE_SECRET_KEY)
	}
	return secretKey
}

type jwtService struct {
	secretKey  string
	issuer     string
	accessExp  time.Duration
	refreshExp time.Duration
}

func NewJwtService() JwtService {
	return &jwtService{
		secretKey:  getSecretKey(),
		issuer:     constants.JWT_ISSUER,
		accessExp:  time.Duration(constants.JWT_ACCESS_EXP),
		refreshExp: time.Duration(constants.JWT_REFRESH_EXP),
	}
}

func (j *jwtService) GenerateAccessToken(userId string, roleName string, userEmail string) (string, error) {
	claims := jwtCustomClaim{
		UserID:    userId,
		RoleName:  roleName,
		UserEmail: userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessExp)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtService) GenerateRefreshToken() (string, time.Time) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Println(err)
		return "", time.Time{}
	}

	refreshToken := base64.StdEncoding.EncodeToString(b)
	expiresAt := time.Now().Add(j.refreshExp)

	return refreshToken, expiresAt
}

func (j *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &jwtCustomClaim{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}

func (j *jwtService) GetUserIDByToken(tokenString string) (string, error) {
	token, err := j.ValidateToken(tokenString)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", fmt.Errorf("token has expired")
		}
		return "", err
	}

	if claims, ok := token.Claims.(*jwtCustomClaim); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", fmt.Errorf("invalid token claims")
}

func (j *jwtService) GetRoleNameByToken(tokenString string) (string, error) {
	token, err := j.ValidateToken(tokenString)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", fmt.Errorf("token has expired")
		}
		return "", err
	}

	if claims, ok := token.Claims.(*jwtCustomClaim); ok && token.Valid {
		return claims.RoleName, nil
	}

	return "", fmt.Errorf("invalid token claims")
}

func (j *jwtService) GetUserEmailByToken(tokenString string) (string, error) {
	token, err := j.ValidateToken(tokenString)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", fmt.Errorf("token has expired")
		}
		return "", err
	}

	if claims, ok := token.Claims.(*jwtCustomClaim); ok && token.Valid {
		return claims.UserEmail, nil
	}

	return "", fmt.Errorf("invalid token claims")
}
