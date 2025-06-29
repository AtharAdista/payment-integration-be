package xenditpay

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"payment/internal/model"

	"github.com/google/uuid"
	xendit "github.com/xendit/xendit-go/v7"
	"github.com/xendit/xendit-go/v7/balance_and_transaction"
	"github.com/xendit/xendit-go/v7/customer"
	"github.com/xendit/xendit-go/v7/invoice"
	"github.com/xendit/xendit-go/v7/payment_method"
	payment_request "github.com/xendit/xendit-go/v7/payment_request"
)

func GetBalance() (*balance_and_transaction.Balance, error) {
	client := xendit.NewClient(os.Getenv("XENDIT_API_KEY"))

	accountType := "CASH"
	currency := "IDR"

	resp, _, err := client.BalanceApi.GetBalance(context.Background()).
		AccountType(accountType).
		Currency(currency).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	return resp, nil
}

func CreatePayment(req *model.SubscriptionPaymentReq) (*payment_request.PaymentRequest, error) {

	xenditClient := xendit.NewClient(os.Getenv("XENDIT_API_KEY"))

	paymentRequestParameters := *payment_request.NewPaymentRequestParameters(payment_request.PAYMENTREQUESTCURRENCY_IDR)

	paymentRequestParameters.SetAmount(req.PaymentRequest.Amount)
	paymentRequestParameters.SetReferenceId(req.PaymentRequest.ReferenceID)
	paymentRequestParameters.SetDescription(req.PaymentRequest.Description)
	paymentRequestParameters.SetCustomerId(req.PaymentRequest.CustomerID)

	paymentMethod := &payment_request.PaymentMethodParameters{
		Type:        payment_request.PaymentMethodType(req.PaymentRequest.PaymentMethod.Type),
		Reusability: payment_request.PaymentMethodReusability(req.PaymentRequest.PaymentMethod.Reusability),
	}

	switch req.PaymentRequest.PaymentMethod.Type {
	case "QR_CODE":
		if req.PaymentRequest.PaymentMethod.QRCode != nil {
			channel := payment_request.QRCodeChannelCode(req.PaymentRequest.PaymentMethod.QRCode.ChannelCode)

			qrParams := payment_request.QRCodeParameters{
				ChannelCode: *payment_request.NewNullableQRCodeChannelCode(&channel),
				ChannelProperties: &payment_request.QRCodeChannelProperties{
					ExpiresAt: req.PaymentRequest.ExpiresAt,
				},
			}
			paymentMethod.QrCode = *payment_request.NewNullableQRCodeParameters(&qrParams)

		}
	case "VIRTUAL_ACCOUNT":
		if req.PaymentRequest.PaymentMethod.VirtualAccount != nil {
			channel := payment_request.VirtualAccountChannelCode(req.PaymentRequest.PaymentMethod.VirtualAccount.ChannelCode)

			vaParams := payment_request.VirtualAccountParameters{
				ChannelCode: payment_request.VirtualAccountChannelCode(channel),
				ChannelProperties: payment_request.VirtualAccountChannelProperties{
					CustomerName: req.PaymentRequest.PaymentMethod.VirtualAccount.ChannelProperties.CustomerName,
					ExpiresAt:    req.PaymentRequest.ExpiresAt,
				},
			}
			paymentMethod.VirtualAccount = *payment_request.NewNullableVirtualAccountParameters(&vaParams)
		}
	default:
		return nil, fmt.Errorf("unsupported payment method type")
	}

	paymentRequestParameters.SetPaymentMethod(*paymentMethod)

	resp, _, err := xenditClient.PaymentRequestApi.CreatePaymentRequest(context.Background()).
		PaymentRequestParameters(paymentRequestParameters).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	return resp, nil

}

func CreateInvoice(req *model.SubscriptionPaymentReq) (*invoice.Invoice, error) {
	xenditClient := xendit.NewClient(os.Getenv("XENDIT_API_KEY"))

	createInvoiceRequest := *invoice.NewCreateInvoiceRequest(req.PaymentRequest.ReferenceID, req.PaymentRequest.Amount)

	createInvoiceRequest.Customer = &invoice.CustomerObject{}
	createInvoiceRequest.Customer.SetCustomerId(req.PaymentRequest.CustomerID)
	createInvoiceRequest.SetPaymentMethods([]string{"CREDIT_CARD", "DEBIT_CARD"})
	createInvoiceRequest.SetInvoiceDuration(3600)
	createInvoiceRequest.SetDescription(req.PaymentRequest.Description)
	createInvoiceRequest.SetSuccessRedirectUrl("http://localhost:5173/")

	resp, _, err := xenditClient.InvoiceApi.CreateInvoice(context.Background()).
		CreateInvoiceRequest(createInvoiceRequest).
		Execute()

	fmt.Println(resp.GetId())

	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	return resp, err

}

func CreateCustomer(req *CreateCustomerInput) (*customer.Customer, *http.Response, error) {

	xenditClient := xendit.NewClient(os.Getenv("XENDIT_API_KEY"))

	referenceId := uuid.NewString()

	individualDetail := customer.IndividualDetail{
		GivenNames: &req.Name,
	}

	customerRequest := *customer.NewCustomerRequest(referenceId)
	customerRequest.SetType("INDIVIDUAL")
	customerRequest.SetEmail(req.Email)
	customerRequest.SetIndividualDetail(individualDetail)

	resp, r, err := xenditClient.CustomerApi.CreateCustomer(context.Background()).
		CustomerRequest(customerRequest).
		Execute()

	if err != nil {
		return nil, r, fmt.Errorf("error: %w", err)
	}

	return resp, r, nil
}

func GetPaymentRequestById(paymentRequestId string) (*payment_request.PaymentRequest, *http.Response, error) {
	xenditClient := xendit.NewClient(os.Getenv("XENDIT_API_KEY"))

	resp, r, err := xenditClient.PaymentRequestApi.GetPaymentRequestByID(context.Background(), paymentRequestId).
		Execute()

	if err != nil {
		return nil, r, fmt.Errorf("error: %w", err)
	}

	return resp, r, nil
}

func GetInvoicePaymentRequestById(paymentRequestId string) (*invoice.Invoice, *http.Response, error) {
	xenditClient := xendit.NewClient(os.Getenv("XENDIT_API_KEY"))

	resp, r, err := xenditClient.InvoiceApi.GetInvoiceById(context.Background(), paymentRequestId).
        Execute()

	if err != nil {
		return nil, r, fmt.Errorf("error: %w", err)
	}

	return resp, r, nil
}

func SimulatePayment(paymentMethodId string, amount float64) error {
	xenditClient := xendit.NewClient(os.Getenv("XENDIT_API_KEY"))

	simulatePaymentRequest := *payment_method.NewSimulatePaymentRequest()
	simulatePaymentRequest.SetAmount(amount)

	fmt.Print(amount)

	_, err := xenditClient.PaymentMethodApi.SimulatePayment(context.Background(), paymentMethodId).
		SimulatePaymentRequest(simulatePaymentRequest).
		Execute()

	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	return nil
}
