package repository

import (
	"database/sql"
	"fmt"
	"payment/internal/errors"
	"payment/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type AuthenticationRepository struct {
	db *sql.DB
}

func NewAuthenticationRepository(db *sql.DB) *AuthenticationRepository {
	return &AuthenticationRepository{db: db}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (r *AuthenticationRepository) CreateUser(user *model.RegisterUserReq) (string, error) {

	tx, err := r.db.Begin()

	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	if err != nil {
		return "", err
	}

	var userAccountID int

	err = tx.QueryRow(`
		INSERT INTO user_account (email, name, password, role_id, customer_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, user.Email, user.Name, user.Password, user.RoleID, user.CustomerId).Scan(&userAccountID)

	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return string(userAccountID), nil

}

func (r *AuthenticationRepository) CheckEmailExist(email string) (bool, error) {
	var user model.UserAccount

	err := r.db.QueryRow(
		`SELECT id from user_account where email=$1`, email,
	).Scan(&user.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, errors.ErrEmailAlreadyExist
}

func (r *AuthenticationRepository) FindUserByEmail(email string) (*model.UserAccount, error) {
	var user model.UserAccount

	err := r.db.QueryRow(
		`SELECT id, email, password, role_id from user_account where email=$1`,
		email,
	).Scan(&user.ID, &user.Email, &user.Password, &user.RoleID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrEmailOrPassWordFalse
		}
		return nil, err
	}

	return &user, nil
}
