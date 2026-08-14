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

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error{
	query := `INSERT INTO users (telegram_id, username, first_name, last_name is_admin, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	return r.db.QueryRow(ctx, query,
		user.TelegramID,
		user.UserName,
		user.FirstName,
		user.LastName,
		user.IsAdmin,
		time.Now(),
	).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepo) GetUser(ctx context.Context, telegramID int64) (*domain.User, error) {
	query := `SELECT id, telegram_id, username, first_name, last_name, is_admin, created_at
		FROM users
		WHERE telegram_id = $1`

	var user domain.User

	err := r.db.QueryRow(ctx, query, telegramID).Scan(
		&user.ID,
		&user.TelegramID,
		&user.UserName,
		&user.FirstName,
		&user.IsAdmin,
		&user.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		fmt.Printf("User with this telegram ID - %d not found", telegramID)
		return nil, nil
	}

	if err != nil {
		fmt.Printf("Unknown error in creating user")
		return nil, err
	}

	return &user, nil
}
