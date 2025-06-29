package handler

import (
	"errors"
	"fmt"
	"net/http"
	customError "payment/internal/errors"
	"payment/internal/model"
	"payment/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthenticationHandler struct {
	authenticationService *service.AuthenticationService
}

func NewAuthenticationHandler(authenticationService *service.AuthenticationService) *AuthenticationHandler {
	return &AuthenticationHandler{authenticationService: authenticationService}
}

func (h *AuthenticationHandler) Register(c *gin.Context) {
	var user model.RegisterUserReq

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	_, err := h.authenticationService.Register(&user)

	if err != nil {
		if errors.Is(err, customError.ErrEmailAlreadyExist) {
			c.JSON(http.StatusConflict, gin.H{"error": customError.ErrEmailAlreadyExist})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create user: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"email": user.Email,
		}})
}

func (h *AuthenticationHandler) Login(c *gin.Context) {
	var user model.LoginUserReq

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result, accessToken, err := h.authenticationService.Login(user.Email, user.Password)
	if err != nil {

		if errors.Is(err, customError.ErrEmailOrPassWordFalse) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": customError.ErrEmailOrPassWordFalse})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"email":        result.User.Email,
			"access_token": accessToken,
		}})
}

func (h *AuthenticationHandler) GetDataUser(c *gin.Context) {
	user, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, user)
}
