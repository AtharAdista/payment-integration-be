package utils

import (
	"fmt"
	"payment/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
	return &JWTMaker{secretKey: secretKey}
}

func (maker *JWTMaker) CreateToken(id int, email string, role int8, duration time.Duration) (string, *model.UserClaims, error) {
	claims, err := NewUserClaims(id, email, role, duration)

	if err != nil {
		return "", nil, fmt.Errorf("error: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(maker.secretKey))

	if err != nil {
		return "", nil, fmt.Errorf("error signing token: %w", err)
	}

	return tokenStr, claims, nil
}

func (maker *JWTMaker) VerifyToken(tokenStr string) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}

		return []byte(maker.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*model.UserClaims)

	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
