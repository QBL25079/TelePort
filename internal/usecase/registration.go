package usecase

import (
	"context"
	"fmt"

	"github.com/QBL25079/TelePort/internal/domain"
)

type Registration struct {
	users domain.UserRepository
}

func NewRegistration(user domain.UserRepository) *Registration {
	return &Registration{users: user}
}

func (uc *Registration) RegisterOrGet(ctx context.Context, telegramID int64, username, firstName string) (*domain.User, error) {
	user, err := uc.users.GetUser(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	// Пользователь уже есть
	if user != nil {
		return user, nil
	}

	// Создаём нового
	user = &domain.User{
		TelegramID: telegramID,
		UserName:   username,
		FirstName:  firstName,
		IsAdmin:    false,
	}

	if err := uc.users.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}