package main

import (
	"context"

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
	subRepo := postgres.NewSubscriptions(db)
	stateRepo := postgres.NewStateRepo(db)

	registration := usecase.NewRegistration(userRepo)
	subscription := usecase.NewSubscription(subRepo, userRepo)
	bot, err := telegram.NewBot(cfg, log, registration, stateRepo, subscription)
	if err != nil {
		log.Fatal("failed to create bot: %w", zap.Error(err))
	}
	bot.Setup()

}
