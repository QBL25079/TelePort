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
	gate  domain.GateClient
}

func NewSubscription(subs domain.SubscriptionRepository, users domain.UserRepository, gate domain.GateClient) *Subscription {
	return &Subscription{subs: subs, users: users, gate: gate}
}

func (s *Subscription) ActivateForUser(ctx context.Context, userID int64, telegramID int64, planID string, locations []string) (*domain.Subscription, error) {
	days, ok := planDays[planID]
	if !ok {
		return nil, fmt.Errorf("unknown plan: %s", planID)
	}

	now := time.Now()
	exp := now.AddDate(0, 0, days)

	link, err := s.gate.CreateOrUpdate(ctx, domain.GateWay{
		TelegramID: telegramID,
		UserName:   fmt.Sprintf("tg_%d", telegramID),
		ExpiresAt:  exp,
	})
	if err != nil {
		return nil, fmt.Errorf("gate: %w", err)
	}

	sub := &domain.Subscription{
		UserID:      userID,
		PlanID:      planID,
		Status:      domain.StatusActive,
		HappLink:    link,
		StartsAt:    now,
		ExpiresAt:   exp,
		LocationIDs: locations,
	}

	if err := s.subs.Create(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *Subscription) GetActive(ctx context.Context, telegramID int64) (*domain.Subscription, error) {
	user, err := s.users.GetUser(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return nil, nil
	}

	return s.subs.GetActiveByUserID(ctx, user.ID)
}