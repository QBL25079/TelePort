package domain

import (
	"context"
	"time"
)

type SubscriptionStatus string

const (
	StatusActive    SubscriptionStatus = "active"
	StatusExpired   SubscriptionStatus = "expired"
	StatusCancelled SubscriptionStatus = "cancelled"
)

type Subscription struct {
	ID        int64
	UserID    int64
	PlanID    string
	Status    SubscriptionStatus
	HappLink  string
	StartsAt  time.Time
	ExpiresAt time.Time
	CreatedAt time.Time

	LocationIDs []string
}

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *Subscription) error
	GetActiveByUserID(ctx context.Context, userID int64) (*Subscription, error)
}
