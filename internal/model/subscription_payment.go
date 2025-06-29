package model

import (
	"time"
)

type Subscription struct {
	UserID    int       `json:"user_id"`
	PackageID int       `json:"package_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type SubscriptionPaymentReq struct {
	UserID         int       `json:"user_id"`
	PackageID      int       `json:"package_id"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	PaymentRequest Payment   `json:"payment"`
}

type Payment struct {
	SubscriptionID  int           `json:"subscription_id"`
	ReferenceID     string        `json:"reference_id"`
	XenditPaymentID string        `json:"xendit_payment_id"`
	CustomerID      string        `json:"customer_id"`
	Description     string        `json:"description"`
	Amount          float64       `json:"amount"`
	Status          string        `json:"status"`
	PaidAt          *time.Time    `json:"paid_at"`
	CreatedAt       time.Time     `json:"created_at"`
	ExpiresAt       *time.Time    `json:"expires_at"`
	PaymentMethod   PaymentMethod `json:"payment_method,omitempty"`
}

type PaymentMethod struct {
	Type           string        `json:"type"`
	Reusability    string        `json:"reusability,omitempty"`
	QRCode         *QRCodeParams `json:"qr_code,omitempty"`
	VirtualAccount *VAParams     `json:"virtual_account,omitempty"`
}

type QRCodeParams struct {
	ChannelCode string `json:"channel_code"`
}

type VAParams struct {
	ChannelCode       string         `json:"channel_code"`
	ChannelProperties VAChannelProps `json:"channel_properties"`
}

type VAChannelProps struct {
	CustomerName string `json:"customer_name"`
}

type GetPaymentRequestByIdRes struct {
	PaymentMethod struct {
		ID             string `json:"id"`
		VirtualAccount *struct {
			ChannelProperties struct {
				CustomerName         string `json:"customer_name"`
				VirtualAccountNumber string `json:"virtual_account_number"`
			} `json:"channel_properties"`
		} `json:"virtual_account,omitempty"`
	} `json:"payment_method"`
}
