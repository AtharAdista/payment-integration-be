package model

import "time"

type RegisterUserReq struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	CustomerId string `json:"customer_id"`
	Password   string `json:"password"`
	RoleID     int    `json:"role_id"`
}

type RegisterUserRes struct {
	ID     int    `json:"id"`
	Email  string `json:"email"`
	RoleID int    `json:"role_id"`
}

type LoginUserRes struct {
	AccessToken          string          `json:"access_token"`
	AccessTokenExpiresAt time.Time       `json:"access_token_expires_at"`
	User                 RegisterUserRes `json:"user"`
}

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
