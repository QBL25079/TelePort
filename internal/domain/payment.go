package domain

import "time"

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "pending"
	PaymentSuccess PaymentStatus = "success"
	PaymentFailed  PaymentStatus = "failed"
)

type Payment struct {
	ID         int
	UserID     int64
	PlanID     string
	Amount     int
	Currency   string
	Status     PaymentStatus
	ExternalID string
	CreatedAt  time.Time
}

type PaymentRepository interface {
	
}
