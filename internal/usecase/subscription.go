package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/QBL25079/TelePort/internal/domain"
)

var planDays = map[string]int{
	"1m": 30, "3m": 90, "6m": 180, "12m": 365,
}

type Subscription struct {
	subs  domain.SubscriptionRepository
	users domain.UserRepository
}

func NewSubscription(subs domain.SubscriptionRepository, users domain.UserRepository) *Subscription {
	return &Subscription{subs: subs, users: users}
}

func (s *Subscription) ActivateFromState(ctx context.Context, telegramID int64, planID string, locations []string) (*domain.Subscription, error) {
	user, err := s.users.GetUser(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("User with this telegram id not found: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	days, ok := planDays[planID]
	if !ok {
		return nil, fmt.Errorf("unknown plan: %s", planID)
	}

	now := time.Now()
	sub := &domain.Subscription{
		UserID:      user.ID,
		PlanID:      planID,
		Status:      domain.StatusActive,
		HappLink:    "",
		StartsAt:    now,
		ExpiresAt:   now.AddDate(0, 0, days),
		LocationIDs: locations}

	if err := s.subs.Create(ctx, sub); err != nil {
		return nil, fmt.Errorf("create subscription: %w", err)
	}
	return sub, nil
}

func (s *Subscription) GetActive(ctx context.Context, telegramID int64) (*domain.Subscription, error) {
	user, err := s.users.GetUser(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("Error to get user: %w", err)
	}

	return s.subs.GetActiveByUserID(ctx, user.ID)
}

func (s *Subscription) ActivateForUser(ctx context.Context, userID int64, planID string, locations []string) (*domain.Subscription, error) {
	days, ok := planDays[planID]
	if !ok {
		return nil, fmt.Errorf("unknown plan: %s", planID)
	}

	now := time.Now()

	sub := &domain.Subscription{
		UserID:      userID,
		PlanID:      planID,
		Status:      domain.StatusActive,
		HappLink:    "",
		StartsAt:    now,
		ExpiresAt:   now.AddDate(0, 0, days),
		LocationIDs: locations,
	}

	if err := s.subs.Create(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}
