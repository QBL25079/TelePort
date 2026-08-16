package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/QBL25079/TelePort/internal/domain"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepo struct {
	db *pgxpool.Pool
}

func NewPaymentRepo(db *pgxpool.Pool) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) Create(ctx context.Context, p *domain.Payment) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO payments (user_id, plan_id, amount, currency, status, locations, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`,
		p.UserID, p.PlanID, p.Amount, p.Currency, p.Status, p.Locations, time.Now(),
	).Scan(&p.ID, &p.CreatedAt)
}

func (r *PaymentRepo) GetByID(ctx context.Context, id int64) (domain.Payment, error) {
	var payment domain.Payment
	var paidAt *time.Time

	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, plan_id, amount, currency, status,
		COALESCE(external_id, ''), locations, created_at, paid_at
		FROM payments WHERE id = $1`).Scan(&payment.ID, &payment.UserID, &payment.PlanID, &payment.Amount, &payment.Currency, &payment.Status,
		&payment.ExternalID, &payment.Locations, &payment.CreatedAt, &paidAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Payment{}, fmt.Errorf("No payments with this ID: %w", err)
		}
		return domain.Payment{}, err
	}

	payment.PaidAt = paidAt
	return payment, nil
}

func (r *PaymentRepo) MarkSuccess(ctx context.Context, id int64) error {
	ct, err := r.db.Exec(ctx, `
		UPDATE payments
		SET status = 'success', paid_at = NOW()
		WHERE id = $1 AND status = 'pending'
	`, id)
	if err != nil {
		return fmt.Errorf("Error to update status")
	}

	if ct.RowsAffected() == 0 {
		return errors.New("payment not pending or not found")
	}
	return nil
}
