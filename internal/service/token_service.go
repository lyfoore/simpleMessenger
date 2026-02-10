package service

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	AccessTokenExpireDuration = 3600
)

type jwtService struct {
	secretKey string
}

type TokenService interface {
	GenerateToken(userID uint) (string, error)
	VerifyToken(tokenString string) (uint, error)
}

func NewJwtService(secretKey string) TokenService {
	return &jwtService{secretKey: secretKey}
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *jwtService) GenerateToken(userID uint) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * AccessTokenExpireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "user_token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) VerifyToken(tokenString string) (uint, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}
	return claims.UserID, nil
}
