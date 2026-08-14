package domain

import (
	"context"
	"time"

)

type User struct {
	ID         int64
	TelegramID int64
	UserName   string
	FirstName  string
	LastName   *string
	IsAdmin    bool
	CreatedAt  time.Time
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetUser(ctx context.Context, telegramID int64) (*User, error)
}