package callback

import tele "gopkg.in/telebot.v4"

func (h *Handler) Pay(c tele.Context) error {
	return c.Edit("Оплата скоро будет подключена.\nПока это заглушка.")
}
