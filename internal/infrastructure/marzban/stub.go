package marzban

import (
	"context"
	"fmt"

	"github.com/QBL25079/TelePort/internal/domain"
)

type Stub struct{}

func NewStub() *Stub {
	return &Stub{}
}

func (s *Stub) CreateOrUpdate(ctx context.Context, u domain.GateWay) (string, error) {
	username := u.UserName
	if username == "" {
		username = fmt.Sprintf("tg_%d", u.TelegramID)
	}
	// Заглушка — имитирует ссылку
	return fmt.Sprintf("https://vpn.example.com/sub/%s", username), nil
}

func (s *Stub) Disable(ctx context.Context, userName string) error {
	return nil
}
