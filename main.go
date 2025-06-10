package main

import (
	"fmt"
	"log/slog"
	"os"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/joho/godotenv"

	"github.com/meesooqa/go-tg-bnews/internal/applog"
	"github.com/meesooqa/go-tg-bnews/internal/config"
)

func main() {
	fmt.Printf("Hello World")
	err := run()
	if err != nil {
		slog.Error("Error running application", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	cp := config.NewDefaultConfigProvider()
	conf, err := cp.GetAppConfig()
	if err != nil {
		return fmt.Errorf("error getting app config: %w", err)
	}

	lp := applog.NewConsoleLoggerProvider(conf.Log)
	logger, cleanup := lp.GetLogger()
	defer cleanup()

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	apiEndpoint := os.Getenv("TELEGRAM_API_ENDPOINT")
	var bot *tgbotapi.BotAPI
	if apiEndpoint == "" {
		bot, err = tgbotapi.NewBotAPI(token)
	} else {
		bot, err = tgbotapi.NewBotAPIWithAPIEndpoint(token, apiEndpoint)
	}
	if err != nil {
		return fmt.Errorf("error creating bot API: %w", err)
	}
	// bot.Debug = true

	logger.Info("Authorized", slog.String("Account", bot.Self.UserName))

	return nil
}
