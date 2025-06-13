package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/meesooqa/go-tg-bnews/internal/config"
	"github.com/meesooqa/go-tg-bnews/internal/proc"
	mytg "github.com/meesooqa/go-tg-bnews/internal/telegram"
)

func main() {
	// TODO periodically check for new messages
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	conf, err := config.Load("etc/config.yml")
	if err != nil {
		return fmt.Errorf("error getting app config: %w", err)
	}

	logger, cleanup := getLogger(isTestMode())
	defer cleanup()
	opts := telegram.Options{
		Logger: logger,
	}
	if isTestMode() {
		opts.DC = 2 // See https://my.telegram.org/apps configuration
		opts.DCList = dcs.Test()
	}
	client, err := telegram.ClientFromEnvironment(opts)
	if err != nil {
		return fmt.Errorf("error creating Telegram client: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	flow := mytg.NewTelegramAuthFlow()
	return client.Run(ctx, func(ctx context.Context) error {
		state := proc.NewPipelineState(ctx, conf, client)
		pipeline := proc.Chain(
			proc.AuthProcessor(flow),
			proc.InitServiceProcessor(),
			proc.LoadChannelsProcessor(os.Getenv("CHANNEL_FROM"), os.Getenv("CHANNEL_TO")),
			proc.FetchMessagesProcessor(),
			proc.FilterProcessor(
				proc.SkipNoText,
				proc.KeepLightning,
			),
			proc.ForwardProcessor(),
		)
		return pipeline(state)
	})
}

func getLogger(isDevelopment bool) (*zap.Logger, func()) {
	var cfg zap.Config
	if isDevelopment {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	cfg.Encoding = "json"
	cfg.OutputPaths = []string{"var/log/app.log"}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	cleanup := func() {
		logger.Sync() // nolint:errcheck
	}
	return logger, cleanup
}

func isTestMode() bool {
	testMode := os.Getenv("IS_TEST_MODE")
	if testMode == "" {
		testMode = "false"
	}
	result, err := strconv.ParseBool(testMode)
	if err != nil {
		//return true, fmt.Errorf("invalid boolean value %q: %v", testMode, err)
		return true
	}
	return result
}
