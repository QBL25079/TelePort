package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PurchaseState struct {
	TelegramID int64
	PlanID     string
	Locations  []string
}

type StateRepo struct {
	db *pgxpool.Pool
}

func NewStateRepo(db *pgxpool.Pool) *StateRepo {
	return &StateRepo{db: db}
}

func (r *StateRepo) Set(ctx context.Context, state *PurchaseState) error {
	query  := `INSERT INTO user_states (telegram_id, plan_id, locations, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (telegram_id) DO UPDATE
		SET plan_id = EXCLUDED.plan_id,
		    locations = EXCLUDED.locations,
		    updated_at = EXCLUDED.updated_at`

	_, err := r.db.Exec(ctx, query,
		state.TelegramID,
		state.PlanID,
		state.Locations,
		time.Now(),
	)
	return err
}

func (r *StateRepo) Get(ctx context.Context, telegramID int64) (*PurchaseState, error) {
	query := `SELECT telegram_id, plan_id, locations
		FROM user_states
		WHERE telegram_id = $1`

	var state PurchaseState
	err := r.db.QueryRow(ctx, query, telegramID).Scan(&state.TelegramID, &state.PlanID, &state.Locations)
	if err != nil {
		return nil, fmt.Errorf("Error in scanning state: %w", err)
	}

	return &state, nil
}

func (r *StateRepo) Clear(ctx context.Context, telegramID int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM user_states WHERE telegram_id = $1`, telegramID)
	return err
}