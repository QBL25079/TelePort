package domain

import (
	"context"
	"time"
)

type GateWay struct {
	TelegramID int64
	UserName   string
	ExpiresAt  time.Time
}

type GateClient interface {
	CreateOrUpdate(ctx context.Context, u GateWay) (string, error)
	Disable(ctx context.Context, userName string) error
}