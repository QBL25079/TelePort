package callback

import (
	"context"
	"fmt"
	"strings"

	"github.com/QBL25079/TelePort/vpn-bot/internal/delivery/telegram/view"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (h *Handler) ToggleLocation(c tele.Context) error {
	ctx := context.Background()
	locID := c.Data()

	state, err := h.State.Get(ctx, c.Sender().ID) 
	if err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Сначала выберите тариф"})
	}

	found := false 
	newLocs := make([]string, 0, len(state.Locations))

	for _, i := range state.Locations {
		if i == locID {
			found = true
			continue
		} 
		newLocs = append(newLocs, i)
	}

	if !found {
		newLocs = append(newLocs, locID)
	}

	state.Locations = newLocs

	if err := h.State.Set(ctx, state); err != nil {
		h.Log.Error("update locations", zap.Error(err))
		return c.Respond(&tele.CallbackResponse{Text: "Ошибка"})
	}

	selected := map[string]bool{}
	for _, l := range state.Locations {
		selected[l] = true
	}

	return c.Edit("Выберите страны:", view.LocationsKeyboard(selected))
}

func (h *Handler) LocationsDone(c tele.Context) error {
	ctx := context.Background()

	state, err := h.State.Get(ctx, c.Sender().ID)
	if err != nil || state == nil || len(state.Locations) == 0 {
		return c.Respond(&tele.CallbackResponse{Text: "Выберите хотя бы одну страну"})
	}

	text := fmt.Sprintf(
		"Проверьте заказ:\n\nТариф: <b>%s</b>\nСтраны: <b>%s</b>\n\nДальше будет оплата.",
		state.PlanID,
		strings.Join(state.Locations, ", "),
	)

	return c.Edit(text, view.ConfirmKeyboard())
}

func (h *Handler) BackToLocations(c tele.Context) error {
	ctx := context.Background()
	state, err := h.State.Get(ctx, c.Sender().ID)
	if err != nil || state == nil {
		return c.Edit("Выберите срок подписки:", view.PlansKeyboard())
	}

	selected := map[string]bool{}
	for _, l := range state.Locations {
		selected[l] = true
	}
	return c.Edit("Выберите страны:", view.LocationsKeyboard(selected))
} 