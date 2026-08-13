package command

import (
	"context"
	"fmt"

	"github.com/QBL25079/TelePort/vpn-bot/internal/delivery/telegram/view"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (h *Handler) Start(c tele.Context) error {
	ctx := context.Background()
	user, err := h.Registration.RegisterOrGet(ctx, c.Sender().ID, c.Sender().Username, c.Sender().FirstName)

	if err != nil {
		h.log.Error("register failed", zap.Error(err))
		return fmt.Errorf("Error to get user. Please register")
	}

	name := user.FirstName
	if name == "" {
		name = "друг"
	}
	return c.Send("Привет, "+name+"!\nВыбери действие:", view.MainMenu())
}
