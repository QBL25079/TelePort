package telegram

import (
	"time"

	"github.com/QBL25079/TelePort/vpn-bot/internal/config"
	"github.com/QBL25079/TelePort/vpn-bot/internal/lib/logger"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func Logger(log *logger.Logger) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			start := time.Now()
			err := next(c)

			log.Info("telegram update",
				zap.Int64("telegram_id", c.Sender().ID),
				zap.String("username", c.Sender().Username),
				zap.String("text", c.Text()),
				zap.Duration("latency", time.Since(start)),
				zap.Error(err),
			)

			return err
		}
	} 
}

func Recover(log *logger.Logger) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			defer func() {
				if r := recover(); r != nil {
					log.Error("panic recovered", zap.Any("panic", r))
					_ = c.Send("Произошла внутренняя ошибка. Попробуй позже.")
				}
			}()
			return next(c)
		}
	}
}


func AdminOnly(cfg *config.Config) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if !cfg.IsAdmin(c.Sender().ID) {
				return c.Send("Accses denied")
			}
			return next(c)
		}
	}
}