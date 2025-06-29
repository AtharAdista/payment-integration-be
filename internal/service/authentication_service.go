package service

import (
	"fmt"
	"payment/internal/errors"
	"payment/internal/model"
	"payment/internal/repository"
	"payment/internal/utils"
	"payment/internal/xenditpay"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct {
	authenticationRepository *repository.AuthenticationRepository
	tokenMaker               *utils.JWTMaker
}

func NewAuthenticationService(repo *repository.AuthenticationRepository, tokenMaker *utils.JWTMaker) *AuthenticationService {
	return &AuthenticationService{authenticationRepository: repo, tokenMaker: tokenMaker}
}

func (s *AuthenticationService) Register(user *model.RegisterUserReq) (string, error) {

	isExist, err := s.authenticationRepository.CheckEmailExist(user.Email)

	if err != nil {
		return "", fmt.Errorf("failed to CheckEmailExist %w", err)
	}

	if isExist {
		return "", fmt.Errorf("cannot create user: %w", errors.ErrEmailAlreadyExist)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("failed to hash password %w", err)
	}

	user.Password = string(hashedPassword)

	createdCustomer, _, err := xenditpay.CreateCustomer(&xenditpay.CreateCustomerInput{
		Email: user.Email,
		Name:  user.Name,
	})

	if err != nil {
		return "", fmt.Errorf("failed to create Xendit customer: %w", err)
	}

	user.CustomerId = createdCustomer.GetId()

	return s.authenticationRepository.CreateUser(user)
}

func (s *AuthenticationService) Login(email string, password string) (*model.LoginUserRes, string, error) {

	user, err := s.authenticationRepository.FindUserByEmail(email)

	if err != nil {
		return nil, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return nil, "", fmt.Errorf("%w", errors.ErrEmailOrPassWordFalse)
	}

	accessToken, accessClaims, err := s.tokenMaker.CreateToken(user.ID, user.Email, int8(user.RoleID), time.Hour*48)

	if err != nil {
		return nil, "", fmt.Errorf("error creating token: %w", err)
	}

	res := &model.LoginUserRes{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
		User: model.RegisterUserRes{
			ID:     user.ID,
			Email:  user.Email,
			RoleID: user.RoleID,
		},
	}

	return res, accessToken, nil
}
