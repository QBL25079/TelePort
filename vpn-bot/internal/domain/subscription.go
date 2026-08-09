package domain

import "time"

type SubscriptionStatus string

const (
	StatusActive    SubscriptionStatus = "active"
	StatusExpired   SubscriptionStatus = "expired"
	StatusCancelled SubscriptionStatus = "cancelled"
)

type Subsciption struct {
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
}
