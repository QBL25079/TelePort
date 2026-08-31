package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/QBL25079/TelePort/internal/domain"
)

type Activator struct {
	subs domain.SubscriptionRepository
	user domain.UserRepository
	gate domain.GateClient
}

func (a *Activator) ActivateForUser(ctx context.Context, userID int64, telegramID int64, planID string, locations []string) (*domain.Subscription, error) {
	days, ok := planDays[planID]
	if !ok {
		return nil, fmt.Errorf("Unknown plan")
	}

	now := time.Now()
	exp := now.AddDate(0, 0, days)

	link, err := a.gate.CreateOrUpdate(ctx, domain.GateWay{TelegramID: telegramID,
		UserName:  fmt.Sprintf("tg_%d", telegramID),
		ExpiresAt: exp})

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

	if err := a.subs.Create(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}
