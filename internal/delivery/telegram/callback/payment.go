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
	if err != nil || state == nil {
		return c.Edit("Сессия покупки истекла. Начни с «Купить подписку».")
	}

	p, err := h.Payment.CreatePending(ctx, c.Sender().ID, state.PlanID, state.Locations)
	if err != nil {
		h.Log.Error("create payment", zap.Error(err))
		return c.Edit("Не удалось создать счёт. Попробуй позже.")
	}

	_ = h.State.Clear(ctx, c.Sender().ID)

	text := fmt.Sprintf(
		"Счёт <b>#%d</b> создан.\n\nТариф: %s\nСумма: %d ₽\nСтраны: %s\n\n"+
			"Статус: <b>ожидает оплаты</b>\n\n"+
			"Сейчас оплата вручную: напиши в поддержку и укажи номер счёта.\n"+
			"После подтверждения подписка активируется.",
		p.ID, p.PlanID, p.Amount, strings.Join(p.Locations, ", "),
	)
	return c.Edit(text)
}
