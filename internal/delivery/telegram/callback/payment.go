package callback

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (h *Handler) Pay(c tele.Context) error {
	ctx := context.Background()

	state, err := h.State.Get(ctx, c.Sender().ID)
	if err != nil {
		return c.Edit("Сессия покупки истекла. Начни с «Купить подписку».")
	}

	sub, err := h.Subscription.ActivateFromState(
		ctx, 
		c.Sender().ID,
		state.PlanID,
		state.Locations,
	)
	if err != nil {
		h.Log.Error("activate subscription", zap.Error(err))
		return c.Edit("Не удалось активировать подписку. Попробуй позже.")
	}

	_ = h.State.Clear(ctx, c.Sender().ID)

	text := fmt.Sprintf(
		"✅ Подписка активна!\n\nТариф: %s\nДо: %s\nСтраны: %s\n\nСсылка для Happ появится после подключения панели.",
		sub.PlanID,
		sub.ExpiresAt.Format("02.01.2006"),
		strings.Join(sub.LocationIDs, ", "),
	)
	return c.Edit(text)
}
