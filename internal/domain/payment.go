package domain

import (
	"context"
	"time"
)

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
	Locations  []string
	ExternalID string
	CreatedAt  time.Time
	PaidAt     *time.Time
}

type PaymentRepository interface {
	Create(ctx context.Context, p *Payment) error
	GetByID(ctx context.Context, id int64) (*Payment, error)
	GetPendingByUserID(ctx context.Context, userID int64) (*Payment, error)
	MarkSuccess(ctx context.Context, id int64) error
}
