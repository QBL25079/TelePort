package usecase

import (
	"context"
	"fmt"

	"github.com/QBL25079/TelePort/internal/domain"
)

var planPrices = map[string]int{
	"1m": 199, "3m": 499, "6m": 899, "12m": 1490,
}

type Payment struct {
	payment domain.PaymentRepository
	users domain.UserRepository
	subs *Subscription
}

func NewPayment(payment domain.PaymentRepository, users domain.UserRepository, subs *Subscription) *Payment {
	return &Payment{payment: payment, users: users, subs: subs}
}

func (p *Payment) CreatePending(ctx context.Context, telegramID int64, planID string, locations []string) (*domain.Payment, error) {
	user, err := p.users.GetUser(ctx, telegramID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("User with this telegram id not found: %w", err)
	}

	price, ok := planPrices[planID]
	if !ok {
		return nil, fmt.Errorf("unknown plan: %s", planID)
	}

	if len(locations) == 0 {
		return nil, fmt.Errorf("no locations")
	}

	pay := &domain.Payment{
		UserID:    user.ID,
		PlanID:    planID,
		Amount:    price,
		Currency:  "RUB",
		Status:    domain.PaymentPending,
		Locations: locations,
	}

	if err := p.payment.Create(ctx, pay); err != nil {
		return nil, fmt.Errorf("Error in creating payment: %w", err)
	}

	return pay, nil
}

func (p *Payment) Confirm(ctx context.Context, paymentID int64) (*domain.Subscription, error) {
	pay, err := p.payment.GetByID(ctx, paymentID)
	if err != nil || p == nil {
		return nil, fmt.Errorf("payment not found")
	}

	if pay.Status != domain.PaymentPending {
		return nil, fmt.Errorf("payment status is %s", pay.Status)
	}

	if err := p.payment.MarkSuccess(ctx, paymentID); err != nil {
		return nil, err
	}

	return p.subs.ActivateForUser(ctx, pay.UserID, pay.PlanID, pay.Locations)
}
