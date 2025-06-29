package model

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  int8   `json:"role"`
	jwt.RegisteredClaims
}
