package callback

import (
	"context"

	"github.com/QBL25079/TelePort/internal/delivery/telegram/view"
	"github.com/QBL25079/TelePort/internal/repository/postgres"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (h *Handler) ChoosePlan(c tele.Context) error {
	ctx := context.Background()
	planID := c.Data()

	err := h.State.Set(ctx, &postgres.PurchaseState{
		TelegramID: c.Sender().ID,
		PlanID:     planID,
		Locations:  []string{},
	})

	if err != nil {
		h.Log.Error("save plan state", zap.Error(err))
		return c.Respond(&tele.CallbackResponse{Text: "Ошибка"})
	}

	return c.Edit(
		"Тариф выбран.\nВыберите страны (можно несколько):", view.LocationsKeyboard(nil),
	)
}

func (h *Handler) BackToPlans(c tele.Context) error {
	return c.Edit("Выберите срок подписки:", view.PlansKeyboard())
}
