package telegram

import (
	"context"
	"fmt"
	"time"

	"github.com/QBL25079/TelePort/vpn-bot/internal/config"
	"github.com/QBL25079/TelePort/vpn-bot/internal/lib/logger"
	"github.com/QBL25079/TelePort/vpn-bot/internal/repository/postgres"
	"github.com/QBL25079/TelePort/vpn-bot/internal/usecase"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	env          string
	bot          *tele.Bot
	cfg          *config.Config
	log          *logger.Logger
	registration *usecase.Registration
	state *postgres.StateRepo
	// usecases
}

func NewBot(cfg *config.Config, log *logger.Logger, registration *usecase.Registration, state *postgres.StateRepo) (*Bot, error) {
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
		state: state,
	}, nil
}

func (b *Bot) Setup() {
	b.bot.Use(Logger(b.log))
	b.bot.Use(Recover(b.log))

	btnBuy := tele.Btn{Text: "🛒 Купить подписку"}
	btnMy := tele.Btn{Text: "📁 Моя подписка"}
	btnSupport := tele.Btn{Text: "💬 Поддержка"}

	b.bot.Handle(&btnBuy, b.handleBuy)
	b.bot.Handle(&btnMy, b.handleMySubscription)
	b.bot.Handle(&btnSupport, b.handleSupport)
	b.bot.Handle("/start", b.handleStart)
	b.bot.Handle("/admin", AdminOnly(b.cfg)(b.handleAdmin))
	b.bot.Handle(&tele.Btn{Unique: "plan"}, b.handleChoosePlan)
}

func (b *Bot) Start(ctx context.Context) {
	b.log.Info("bot started", zap.String("username", b.bot.Me.Username))
	b.bot.Start()
}

func (b *Bot) Stop() {
	b.log.Info("bot stoped")
	b.bot.Stop()
}

func (b *Bot) TeleBot() *tele.Bot {
	return b.bot
}

func (b *Bot) handleStart(c tele.Context) error {
	return c.Send("Привет! Я VPN-бот.\nСкоро здесь появится меню.")
}

func (b *Bot) handleAdmin(c tele.Context) error {
	return c.Send("Админ-панель")
}

func (b *Bot) handleBuy(c tele.Context) error {
	return c.Send("Здесь будет покупка...")
}

func (b *Bot) handleMySubscription(c tele.Context) error {
	return c.Send("У тебя пока нет активной подписки.")
}

func (b *Bot) handleSupport(c tele.Context) error {
	return c.Send("По всем вопросам пиши @твой_поддержка")
}

func (b *Bot) plansKeyboard() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	btn1 := menu.Data("1 месяц — 199₽", "plan", "1m")
	btn3 := menu.Data("3 месяца — 499₽", "plan", "3m")
	btn6 := menu.Data("6 месяцев — 899₽", "plan", "6m")
	btn12 := menu.Data("1 год — 1490₽", "plan", "12m")
	btnBack := menu.Data("« Назад", "back_main")

	menu.Inline(
		menu.Row(btn1),
		menu.Row(btn3),
		menu.Row(btn6),
		menu.Row(btn12),
		menu.Row(btnBack),
	)

	return menu
}

func (b *Bot) handleChoosePlan(c tele.Context) error {
	planID := c.Data()
	return c.Edit("Вы выбрали тариф: " + planID + "\nТеперь выберите страны (скоро)")
}
