package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/QBL25079/TelePort/internal/config"
	"github.com/QBL25079/TelePort/internal/delivery/telegram"
	"github.com/QBL25079/TelePort/internal/lib/logger"
	"github.com/QBL25079/TelePort/internal/repository/postgres"
	"github.com/QBL25079/TelePort/internal/usecase"
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
		log.Fatal("failed to create bot", zap.Error(err))
	}
	bot.Setup()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go bot.Start() 

	log.Info("bot is running...")
	<-ctx.Done()
	bot.Stop()
}
