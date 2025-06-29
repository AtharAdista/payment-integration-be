package xenditpay

import (
	"fmt"
	"net/http"
	"payment/internal/model"
	"payment/internal/repository"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	subscriptionPaymentRepository *repository.SubscriptionPaymentRepository
}

func NewPaymentHandler(subscriptionPaymentRepository *repository.SubscriptionPaymentRepository) *PaymentHandler {
	return &PaymentHandler{subscriptionPaymentRepository: subscriptionPaymentRepository}
}

func (h *PaymentHandler) PaymentCallbackHandler(c *gin.Context) {
	var payload PaymentWebhookPaylod

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.subscriptionPaymentRepository.UpdatePaymentStatusByRequestID(payload.Data.PaymentRequestID, payload.Data.Status)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to update payment status: %v", err)})
		return
	}

	c.Status(http.StatusOK)
}

func (h *PaymentHandler) InvoiceCallbackHandler(c *gin.Context) {
	var payload InvoiceWebhookPaylod

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.subscriptionPaymentRepository.UpdateInvoicePaymentStatusByRequestID(payload.Id, payload.Status)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to update payment status: %v", err)})
		return
	}

	c.Status(http.StatusOK)
}

func (h *PaymentHandler) SimulatePaymentHandler(c *gin.Context) {
	var payload SimulatePaymentPaylod

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idStr := c.Param("xenditId")

	err := SimulatePayment(idStr, payload.Amount)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to update payment status: %v", err)})
		return
	}

	c.Status(http.StatusOK)
}

func (h *PaymentHandler) GetPaymentRequestByIdHandler(c *gin.Context) {

	userInterface, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	_, ok := userInterface.(*model.UserClaims)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user context"})
		return
	}

	idStr := c.Param("xenditId")

	paymentType := c.Query("type")

	if paymentType == "CARD" {
		invoiceResult, _, err := GetInvoicePaymentRequestById(idStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch invoice: %v", err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"invoice_id":      invoiceResult.Id,
			"invoice_url":     invoiceResult.InvoiceUrl,
		})
		return
	}

	result, _, err := GetPaymentRequestById(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to update payment status: %v", err)})
		return
	}

	if result.PaymentMethod.Type == "VIRTUAL_ACCOUNT" && result.PaymentMethod.VirtualAccount.IsSet() {
		c.JSON(http.StatusCreated, gin.H{
			"payment_id":             result.PaymentMethod.Id,
			"customer_name":          result.PaymentMethod.VirtualAccount.Get().ChannelProperties.CustomerName,
			"virtual_account_number": result.PaymentMethod.VirtualAccount.Get().ChannelProperties.VirtualAccountNumber,
			"type":                   result.PaymentMethod.Type,
		})
	} else if result.PaymentMethod.Type == "QR_CODE" && result.PaymentMethod.QrCode.IsSet() {
		c.JSON(http.StatusCreated, gin.H{
			"payment_id": result.PaymentMethod.Id,
			"type":       result.PaymentMethod.Type,
			"qr_string":  result.PaymentMethod.QrCode.Get().ChannelProperties.QrString,
		})
	}
}
