package callback

import (
	"github.com/QBL25079/TelePort/internal/lib/logger"
	"github.com/QBL25079/TelePort/internal/repository/postgres"
	"github.com/QBL25079/TelePort/internal/usecase"
	tele "gopkg.in/telebot.v4"
)

type Handler struct {
	State        *postgres.StateRepo
	Subscription *usecase.Subscription
	Log          *logger.Logger
	Bot          *tele.Bot
}

func NewHandler(state *postgres.StateRepo, sub *usecase.Subscription, log *logger.Logger, bot *tele.Bot) *Handler {
	return &Handler{State: state, Subscription: sub, Log: log, Bot: bot}
}
