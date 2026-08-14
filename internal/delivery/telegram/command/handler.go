package command

import (
	"github.com/QBL25079/TelePort/internal/lib/logger"
	"github.com/QBL25079/TelePort/internal/usecase"
	tele "gopkg.in/telebot.v4"
)

type Handler struct {
	Registration *usecase.Registration
	Subscription *usecase.Subscription
	log *logger.Logger
	Bot *tele.Bot
}

func NewHandler(reg *usecase.Registration, log *logger.Logger, bot *tele.Bot) *Handler {
	return &Handler{Registration: reg, log: log, Bot: bot}
}

