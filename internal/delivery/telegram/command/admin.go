package command

import tele "gopkg.in/telebot.v4"

func (h *Handler) Admin(c tele.Context) error {
	return c.Send("Админ-панель")
}