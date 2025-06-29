package service

import (
	"fmt"
	"payment/internal/model"
	"payment/internal/repository"
	"payment/internal/xenditpay"
	"time"

	"github.com/google/uuid"
	"github.com/xendit/xendit-go/v7/invoice"
	"github.com/xendit/xendit-go/v7/payment_request"
)

type SubscriptionPaymentService struct {
	subscriptionPaymentRepository *repository.SubscriptionPaymentRepository
}

func NewSubcriptionService(repo *repository.SubscriptionPaymentRepository) *SubscriptionPaymentService {
	return &SubscriptionPaymentService{subscriptionPaymentRepository: repo}
}

func (s *SubscriptionPaymentService) CreateSubscriptionPayment(subscriptionPayment *model.SubscriptionPaymentReq) (string, error) {

	timeNow := time.Now()
	referenceId := uuid.NewString()
	expiresAt := time.Now().Add(1 * time.Hour).UTC()

	subscriptionPayment.StartDate = timeNow
	subscriptionPayment.EndDate = timeNow.AddDate(0, 0, 30)
	subscriptionPayment.PaymentRequest.ReferenceID = referenceId
	subscriptionPayment.PaymentRequest.ExpiresAt = &expiresAt

	customerID, err := s.subscriptionPaymentRepository.FindCustomerId(subscriptionPayment.UserID)

	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	subscriptionPayment.PaymentRequest.CustomerID = customerID

	switch subscriptionPayment.PaymentRequest.PaymentMethod.Type {
	case "QR_CODE":
		if subscriptionPayment.PaymentRequest.PaymentMethod.QRCode == nil {
			return "", fmt.Errorf("qr_code params is nil")
		}

	case "VIRTUAL_ACCOUNT":
		if subscriptionPayment.PaymentRequest.PaymentMethod.VirtualAccount == nil {
			return "", fmt.Errorf("virtual_account params is nil")
		}
		name, err := s.subscriptionPaymentRepository.FindCustomerName(subscriptionPayment.UserID)
		if err != nil {
			return "", fmt.Errorf("failed to find customer name: %w", err)
		}
		subscriptionPayment.PaymentRequest.PaymentMethod.VirtualAccount.ChannelProperties.CustomerName = name

	case "CARD":

	default:
		return "", fmt.Errorf("unsupported payment method: %s", subscriptionPayment.PaymentRequest.PaymentMethod.Type)
	}

	var resultXendit any

	if subscriptionPayment.PaymentRequest.PaymentMethod.Type == "CARD" {
		resultXendit, err = xenditpay.CreateInvoice(subscriptionPayment)
	} else {
		resultXendit, err = xenditpay.CreatePayment(subscriptionPayment)
	}

	if resultXendit == nil {
		return "", fmt.Errorf("failed to create invoice/payment: %w", err)
	}

	switch v := resultXendit.(type) {
	case *invoice.Invoice:
		subscriptionPayment.PaymentRequest.XenditPaymentID = v.GetId()
	case *payment_request.PaymentRequest:
		subscriptionPayment.PaymentRequest.XenditPaymentID = v.GetId()
	default:
		return "", fmt.Errorf("unknown Xendit response type")
	}

	if err := s.subscriptionPaymentRepository.CreateSubscriptionPayment(subscriptionPayment); err != nil {
		return "", fmt.Errorf("failed to save subscription payment: %w", err)
	}

	return subscriptionPayment.PaymentRequest.XenditPaymentID, nil
}

func (s *SubscriptionPaymentService) CheckPaymentExpired(id int) error {
	result, err := s.subscriptionPaymentRepository.FindCustomerId(id)

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return s.subscriptionPaymentRepository.CheckPaymentExpired(result)
}

func (s *SubscriptionPaymentService) CheckPaymentStatusById(xenditId string) (string, error) {
	result, err := s.subscriptionPaymentRepository.CheckPaymentStatusById(xenditId)

	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	return result, nil

}

func (s *SubscriptionPaymentService) GetUserActivePackageName(userId int) (string, error) {
	packageName, err := s.subscriptionPaymentRepository.GetUserActivePackageName(userId)
	if err != nil {
		return "", fmt.Errorf("service failed to get active package: %w", err)
	}

	return packageName, nil
}

func (s *SubscriptionPaymentService) CheckSubscriptionActive(userId int) (bool, error) {
	return s.subscriptionPaymentRepository.CheckSubscriptionActive(userId)
}
