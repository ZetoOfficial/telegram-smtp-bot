package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ZetoOfficial/telegram-smtp-bot/internal/adapter/smtp"
	"github.com/ZetoOfficial/telegram-smtp-bot/internal/adapter/telegram"
	"github.com/ZetoOfficial/telegram-smtp-bot/internal/config"
	"github.com/ZetoOfficial/telegram-smtp-bot/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	smtpConfig := smtp.SMTPConfig{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUsername,
		Password: cfg.SMTPPassword,
	}
	smtpClient := smtp.NewSMTPSender(smtpConfig)

	botUseCase := usecase.NewBotUseCase(smtpClient)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("Ошибка создания Telegram бота: %v", err)
	}

	bot.Debug = false

	telegramHandler := telegram.NewTelegramHandler(bot, botUseCase)

	log.Printf("Бот Telegram запущен. ID: %d", bot.Self.ID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})

	go func() {
		telegramHandler.HandleUpdates(ctx)
		close(done)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	log.Printf("Получен сигнал %s, инициируется завершение работы...", sig)

	cancel()

	select {
	case <-done:
		log.Println("Обработчик Telegram завершил работу.")
	case <-time.After(10 * time.Second):
		log.Println("Таймаут ожидания завершения работы обработчика Telegram.")
	}

	log.Println("Приложение завершено.")
}
