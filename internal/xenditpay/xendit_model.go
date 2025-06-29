package xenditpay

// import (
// 	"github.com/xendit/xendit-go/v7"
// 	"github.com/xendit/xendit-go/v7/payment_request"
// )

type CreateCustomerInput struct {
	Email string
	Name  string
}

type PaymentWebhookPaylod struct {
	Event string                   `json:"event"`
	Data  DataPaymentWebhookPaylod `json:"data"`
}

type DataPaymentWebhookPaylod struct {
	Status           string `json:"status"`
	PaymentRequestID string `json:"payment_request_id"`
}

type InvoiceWebhookPaylod struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type SimulatePaymentPaylod struct {
	Amount float64 `json:"amount"`
}
