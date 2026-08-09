package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/QBL25079/TelePort/vpn-bot/internal/config"
	"github.com/QBL25079/TelePort/vpn-bot/internal/delivery/telegram"
	"github.com/QBL25079/TelePort/vpn-bot/internal/lib/logger"
	"github.com/QBL25079/TelePort/vpn-bot/internal/repository/postgres"
	"github.com/QBL25079/TelePort/vpn-bot/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()

	logCfg := logger.NewMustConfig()
	log, err := logger.NewLogger(logCfg)
	if err != nil {
		panic(err)
	}

	defer log.Close()

	db, err := pgxpool.New(context.Background(), cfg.Postgres.DSN())
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	userRepo := postgres.NewUserRepo(db)
	registration := usecase.NewRegistration(userRepo)
	state := postgres.NewStateRepo(db)

	bot, err := telegram.NewBot(cfg, log, registration, state)
	bot.Setup()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		bot.Start(ctx)
	}()

	log.Info("bot is running...")
	<-ctx.Done()

	log.Info("shutting down...")
	bot.Stop()
}
