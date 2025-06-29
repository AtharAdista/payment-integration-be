package handler

import (
	"fmt"
	"net/http"
	"payment/internal/model"
	"payment/internal/service"

	"github.com/gin-gonic/gin"
)

type SubscriptionPaymentHandler struct {
	subscriptionPaymentService *service.SubscriptionPaymentService
}

func NewSubscriptionPaymentHandler(subscriptionPaymentService *service.SubscriptionPaymentService) *SubscriptionPaymentHandler {
	return &SubscriptionPaymentHandler{subscriptionPaymentService: subscriptionPaymentService}
}

func (h *SubscriptionPaymentHandler) CreateSubscriptionPayment(c *gin.Context) {

	var req model.SubscriptionPaymentReq

	userInterface, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, ok := userInterface.(*model.UserClaims)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user context"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	req.UserID = user.ID

	idXendit, err := h.subscriptionPaymentService.CreateSubscriptionPayment(&req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create subscription: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": idXendit,
	})

}

func (h *SubscriptionPaymentHandler) CheckPaymentExpired(c *gin.Context) {

	userInterface, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, ok := userInterface.(*model.UserClaims)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user context"})
		return
	}

	err := h.subscriptionPaymentService.CheckPaymentExpired(user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to check payment: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id": user.ID,
	})
}

func (h *SubscriptionPaymentHandler) CheckPaymentStatusById(c *gin.Context) {

	_, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("xenditId")

	result, err := h.subscriptionPaymentService.CheckPaymentStatusById(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to payment status: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": result,
	})
}

func (h *SubscriptionPaymentHandler) GetUserPackageAccess(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, ok := userInterface.(*model.UserClaims)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user context"})
		return
	}

	packageName, err := h.subscriptionPaymentService.GetUserActivePackageName(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user package"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"package": packageName})
}

func (h *SubscriptionPaymentHandler) CheckSubscriptionActive(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, ok := userInterface.(*model.UserClaims)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user context"})
		return
	}

	isActive, err := h.subscriptionPaymentService.CheckSubscriptionActive(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check subscription",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"active": isActive,
	})
}
