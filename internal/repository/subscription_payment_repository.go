package repository

import (
	"database/sql"
	"errors"
	"fmt"
	errorsManual "payment/internal/errors"
	"payment/internal/model"
	"time"
)

type SubscriptionPaymentRepository struct {
	db *sql.DB
}

func NewSubscriptionPaymentRepository(db *sql.DB) *SubscriptionPaymentRepository {
	return &SubscriptionPaymentRepository{db: db}
}

func (r *SubscriptionPaymentRepository) FindCustomerId(id int) (string, error) {

	var customerId string

	err := r.db.QueryRow(
		`SELECT customer_id from user_account where id=$1`,
		id,
	).Scan(&customerId)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errorsManual.ErrIdNotFound
		}
		return "", err
	}

	return customerId, nil
}

func (r *SubscriptionPaymentRepository) FindCustomerName(id int) (string, error) {

	var name string

	err := r.db.QueryRow(
		`SELECT name from user_account where id=$1`,
		id,
	).Scan(&name)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errorsManual.ErrIdNotFound
		}
		return "", err
	}

	return name, nil
}

func (r *SubscriptionPaymentRepository) CreateSubscriptionPayment(subscriptionPayment *model.SubscriptionPaymentReq) error {

	tx, err := r.db.Begin()

	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	var subscriptionID int

	err = tx.QueryRow(`
		INSERT INTO subscriptions (user_id, package_id, start_date, end_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, subscriptionPayment.UserID, subscriptionPayment.PackageID, subscriptionPayment.StartDate, subscriptionPayment.EndDate).Scan(&subscriptionID)

	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	subscriptionPayment.PaymentRequest.SubscriptionID = subscriptionID

	_, err = tx.Exec(`
		INSERT INTO payments (subscription_id, reference_id, xendit_payment_id, channel, amount, expired_at, customer_id, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, subscriptionPayment.PaymentRequest.SubscriptionID, subscriptionPayment.PaymentRequest.ReferenceID, subscriptionPayment.PaymentRequest.XenditPaymentID, subscriptionPayment.PaymentRequest.PaymentMethod.Type, subscriptionPayment.PaymentRequest.Amount, subscriptionPayment.PaymentRequest.ExpiresAt, subscriptionPayment.PaymentRequest.CustomerID, subscriptionPayment.PaymentRequest.Description)

	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *SubscriptionPaymentRepository) CheckPaymentExpired(customerId string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE payments
		SET status = 'EXPIRED'
		WHERE status = 'PENDING'
		AND expired_at < CURRENT_TIMESTAMP
		AND customer_id = $1
	`, customerId)
	if err != nil {
		return fmt.Errorf("failed to update expired payments: %w", err)
	}

	_, err = tx.Exec(`
		UPDATE subscriptions
		SET status = 'CANCELLED'
		WHERE id IN (
			SELECT subscription_id
			FROM payments
			WHERE customer_id = $1 AND status = 'EXPIRED'
		)
		AND status = 'PENDING'
	`, customerId)
	if err != nil {
		return fmt.Errorf("failed to cancel expired subscriptions: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *SubscriptionPaymentRepository) UpdatePaymentStatusByRequestID(paymentRequestID, status string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var subscriptionID int
	err = tx.QueryRow(`
		SELECT subscription_id
		FROM payments
		WHERE xendit_payment_id = $1
	`, paymentRequestID).Scan(&subscriptionID)

	if err != nil {
		return fmt.Errorf("failed to fetch subscription_id: %w", err)
	}

	var paidAt interface{}
	if status == "SUCCEEDED" {
		paidAt = time.Now()
	} else {
		paidAt = nil
	}

	_, err = tx.Exec(`
		UPDATE payments
		SET status = $1,
			paid_at = $2
		WHERE xendit_payment_id = $3
	`, status, paidAt, paymentRequestID)

	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	if status == "SUCCEEDED" {
		_, err = tx.Exec(`
			UPDATE subscriptions
			SET status = 'ACTIVE'
			WHERE id = $1
		`, subscriptionID)

		if err != nil {
			return fmt.Errorf("failed to update subscription status: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *SubscriptionPaymentRepository) UpdateInvoicePaymentStatusByRequestID(paymentRequestID, status string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var subscriptionID int
	err = tx.QueryRow(`
		SELECT subscription_id
		FROM payments
		WHERE xendit_payment_id = $1
	`, paymentRequestID).Scan(&subscriptionID)

	if err != nil {
		return fmt.Errorf("failed to fetch subscription_id: %w", err)
	}

	var statusDatabase string
	
	if status == "PAID"{
		statusDatabase = "SUCCEEDED"
	}

	var paidAt interface{}
	if statusDatabase == "SUCCEEDED" {
		paidAt = time.Now()
	} else {
		paidAt = nil
	}

	_, err = tx.Exec(`
		UPDATE payments
		SET status = $1,
			paid_at = $2
		WHERE xendit_payment_id = $3
	`, statusDatabase, paidAt, paymentRequestID)

	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	if statusDatabase == "SUCCEEDED" {
		_, err = tx.Exec(`
			UPDATE subscriptions
			SET status = 'ACTIVE'
			WHERE id = $1
		`, subscriptionID)

		if err != nil {
			return fmt.Errorf("failed to update subscription status: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *SubscriptionPaymentRepository) CheckPaymentStatusById(xenditId string) (string, error) {
	var status string

	err := r.db.QueryRow(`
		SELECT status
		FROM payments
		WHERE xendit_payment_id = $1
	`, xenditId).Scan(&status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return status, fmt.Errorf("failed to check payment status: %w", err)
	}

	return status, nil
}

func (r *SubscriptionPaymentRepository) GetUserActivePackageName(userId int) (string, error) {
	var packageName string

	err := r.db.QueryRow(`
		SELECT p.name
		FROM subscriptions s
		JOIN packages p ON s.package_id = p.id
		WHERE s.user_id = $1 AND s.status = 'ACTIVE'
		LIMIT 1
	`, userId).Scan(&packageName)

	if err != nil {
		if err == sql.ErrNoRows {
			return "free", nil
		}
		return "", fmt.Errorf("failed to get active package: %w", err)
	}

	return packageName, nil
}

func (r *SubscriptionPaymentRepository) CheckSubscriptionActive(userId int) (bool, error) {
	var exists bool

	err := r.db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM subscriptions 
			WHERE user_id = $1 AND status = 'ACTIVE'
		)
	`, userId).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check active subscription: %w", err)
	}

	return exists, nil
}
