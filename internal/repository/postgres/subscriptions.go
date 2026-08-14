package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/QBL25079/TelePort/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubRepo struct {
	db *pgxpool.Pool
}

func NewSubscriptions(db *pgxpool.Pool) *SubRepo {
	return &SubRepo{db: db}
}

func (r *SubRepo) Create(ctx context.Context, sub *domain.Subscription) error {
	tx, err := r.db.Begin(ctx)

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	query := `INSERT INTO subscriptions (user_id, plan_id, status, happ_link, starts_at, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;`

	err = tx.QueryRow(ctx, query, sub.UserID, sub.PlanID, sub.Status, sub.HappLink, sub.StartsAt, sub.ExpiresAt, time.Now()).Scan(&sub.ID)

	if err != nil {
		return fmt.Errorf("Error to insert sub: %w", err)
	}

	for _, locID := range sub.LocationIDs {
		_, err := tx.Exec(ctx, `INSERT INTO subscription_locations (subscription_id, location_id) VALUES ($1, $2);`, sub.ID, locID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *SubRepo) GetActiveByUserID(ctx context.Context, userID int64) (*domain.Subscription, error) {
	query := `SELECT id, user_id, plan_id, status, COALESCE(happ_link, '') starts_at, expires_at, 
	created_at FROM subscriptions WHERE user_id = $1 AND status = 'active' AND expires_at > NOW()
		ORDER BY expires_at DESC
		LIMIT 1`
	var sub domain.Subscription
	
	err := r.db.QueryRow(ctx, query, userID).Scan(&sub.ID, &sub.UserID, &sub.PlanID, &sub.Status, &sub.HappLink,
		&sub.StartsAt, &sub.ExpiresAt, &sub.CreatedAt,)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("this subscription does not exists: %w", err)
		}
		return nil, err
	}

	rows, err := r.db.Query(ctx, `SELECT location_id FROM subscription_locations WHERE subscription_id = $1`, &sub.ID)

	if err != nil {
		return nil, fmt.Errorf("Cant find locs in sub: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var id string 
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("Error to scan subs ID: %w", err)
		}
		sub.LocationIDs = append(sub.LocationIDs, id)
	}

	return &sub, nil
}