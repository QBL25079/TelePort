package command

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (h *Handler) MySubscription(c tele.Context) error {
	ctx := context.Background()

	sub, err := h.Subscription.GetActive(ctx, c.Sender().ID)
	if err != nil {
		h.log.Error("get subscription", zap.Error(err))
		return c.Send("Ошибка. Попробуй позже.")
	}

	if sub == nil {
		return c.Send("У тебя нет активной подписки.")
	}

	text := fmt.Sprintf(
		"📁 Твоя подписка\n\nТариф: %s\nСтатус: %s\nДо: %s\nСтраны: %s",
		sub.PlanID,
		sub.Status,
		sub.ExpiresAt.Format("02.01.2006 15:04"),
		strings.Join(sub.LocationIDs, ", "),
	)
	return c.Send(text)
}
