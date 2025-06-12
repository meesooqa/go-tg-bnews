package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/dcs"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/meesooqa/go-tg-bnews/internal/config"
	"github.com/meesooqa/go-tg-bnews/internal/proc"
	mytg "github.com/meesooqa/go-tg-bnews/internal/telegram"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
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

	ctx := context.Background()

	logger, cleanup := getLogger(isTestMode())
	defer cleanup()
	opts := telegram.Options{
		Logger: logger,
	}
	if isTestMode() {
		opts.DC = 2
		opts.DCList = dcs.Test()
	}
	client, err := telegram.ClientFromEnvironment(opts)
	if err != nil {
		return fmt.Errorf("error creating Telegram client: %w", err)
	}
	flow := auth.NewFlow(mytg.AuthFlow{}, auth.SendCodeOptions{})
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

func getLogger(isDevelopment bool) (logger *zap.Logger, cleanup func()) {
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
	cleanup = func() {
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
