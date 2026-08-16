package command

import (
	"github.com/QBL25079/TelePort/internal/lib/logger"
	"github.com/QBL25079/TelePort/internal/usecase"
	tele "gopkg.in/telebot.v4"
)

type Handler struct {
	Registration *usecase.Registration
	Subscription *usecase.Subscription
	Payment *usecase.Payment
	log *logger.Logger
	Bot *tele.Bot
}

func NewHandler(reg *usecase.Registration, sub *usecase.Subscription, payment *usecase.Payment, log *logger.Logger, bot *tele.Bot) *Handler {
	return &Handler{Registration: reg, Subscription: sub, Payment: payment, log: log, Bot: bot}
}

