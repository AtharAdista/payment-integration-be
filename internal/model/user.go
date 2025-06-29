package model

import "time"

type UserAccount struct {
	ID           int            `json:"id"`
	CustomerId   string         `json:"customer_id"`
	Email        string         `json:"email"`
	Name         string         `json:"name"`
	Password     string         `json:"password"`
	CreatedAt    time.Time      `json:"created_at"`
	RoleID       int            `json:"role_id"`
	Subscription []Subscription `json:"subscription"`
}

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
