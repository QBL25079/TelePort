package telegram

import (
	"fmt"
	"time"

	"github.com/QBL25079/TelePort/internal/config"
	"github.com/QBL25079/TelePort/internal/delivery/telegram/callback"
	"github.com/QBL25079/TelePort/internal/delivery/telegram/command"
	"github.com/QBL25079/TelePort/internal/delivery/telegram/view"
	"github.com/QBL25079/TelePort/internal/lib/logger"
	"github.com/QBL25079/TelePort/internal/repository/postgres"
	"github.com/QBL25079/TelePort/internal/usecase"
	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	env          string
	bot          *tele.Bot
	cfg          *config.Config
	log          *logger.Logger
	registration *usecase.Registration
	state        *postgres.StateRepo
	sub          *usecase.Subscription
	// usecases
}

func NewBot(cfg *config.Config, log *logger.Logger, registration *usecase.Registration, state *postgres.StateRepo, sub *usecase.Subscription) (*Bot, error) {
	pref := tele.Settings{
		Token:     cfg.Bot.Token,
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: tele.ModeHTML,
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("Error to config a New Bot: %w", err)
	}

	return &Bot{
		env:          cfg.Env,
		bot:          bot,
		cfg:          cfg,
		log:          log,
		registration: registration,
		state:        state,
		sub:          sub,
	}, nil
}

func (b *Bot) Setup() {
	b.bot.Use(Logger(b.log))
	b.bot.Use(Recover(b.log))

	cmd := command.NewHandler(b.registration, b.sub, b.log, b.bot)
	cb := callback.NewHandler(b.state, b.sub, b.log, b.bot)

	b.bot.Handle("/start", cmd.Start)
	b.bot.Handle("/admin", AdminOnly(b.cfg)(cmd.Admin))

	b.bot.Handle(tele.OnText, func(c tele.Context) error {
		switch c.Text() {
		case "🛒 Купить подписку":
			return c.Send("Выберите срок подписки:", view.PlansKeyboard())
		case "📁 Моя подписка":
    		return cmd.MySubscription(c)
		case "💬 Поддержка":
			return c.Send("По вопросам: @support")
		default:
			return nil
		}
	})

	b.bot.Handle(&tele.Btn{Unique: "plan"}, cb.ChoosePlan)
	b.bot.Handle(&tele.Btn{Unique: "loc"}, cb.ToggleLocation)
	b.bot.Handle(&tele.Btn{Unique: "loc_done"}, cb.LocationsDone)
	b.bot.Handle(&tele.Btn{Unique: "pay"}, cb.Pay)
	b.bot.Handle(&tele.Btn{Unique: "back_plans"}, cb.BackToPlans)
	b.bot.Handle(&tele.Btn{Unique: "back_locs"}, cb.BackToLocations)
	b.bot.Handle(&tele.Btn{Unique: "back_main"}, cmd.Start)
}

func (b *Bot) Start() {
	b.log.Info("bot started")
	b.bot.Start()
}

func (b *Bot) Stop() {
	b.log.Info("bot stopped")
	b.bot.Stop()
}
