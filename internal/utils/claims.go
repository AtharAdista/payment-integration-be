package utils

import (
	"fmt"
	"payment/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func NewUserClaims(id int, email string, role int8, duration time.Duration) (*model.UserClaims, error) {
	tokenID, err := uuid.NewRandom()

	if err != nil {
		return nil, fmt.Errorf("error generating token ID: %w", err)
	}

	return &model.UserClaims{
		Email: email,
		ID:    id,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}
