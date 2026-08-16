package command

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (h *Handler) Admin(c tele.Context) error {
	return c.Send("Админ-панель")
}

func (h *Handler) ConfirmPayment(c tele.Context) error {
	args := strings.Fields(c.Text())
	if len(args) < 2 {
		return c.Send("Использование: /confirm_payment <id>")
	}

	id, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return c.Send("Некорректный id")
	}

	sub, err := h.Payment.Confirm(context.Background(), id)
	if err != nil {
		h.log.Error("confirm payment", zap.Error(err))
		return c.Send("Ошибка: " + err.Error())
	}

	return c.Send(fmt.Sprintf("Оплата #%d подтверждена. Подписка до %s", id, sub.ExpiresAt.Format("02.01.2006")))
}
