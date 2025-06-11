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

	"github.com/meesooqa/go-tg-bnews/internal/applog"
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
	lp := applog.NewFileLoggerProvider(conf.Log)
	logger, cleanup := lp.GetLogger()
	defer cleanup()

	ctx := context.Background()

	appID, _ := strconv.Atoi(os.Getenv("APP_ID"))
	appHash := os.Getenv("APP_HASH")
	flow := auth.NewFlow(mytg.AuthFlow{}, auth.SendCodeOptions{})
	client := telegram.NewClient(appID, appHash, telegram.Options{
		DC:     2,
		DCList: dcs.Test(),
	})
	return client.Run(ctx, func(ctx context.Context) error {
		state := &proc.PipelineState{
			Ctx:    ctx,
			Client: client,
		}
		pipeline := proc.Chain(
			proc.AuthProcessor(flow),
			proc.InitServiceProcessor(logger),
			proc.LoadChannelsProcessor("test_bbbolt_001", "test_bbbolt_002"),
			proc.FetchMessagesProcessor(),
			proc.FilterProcessor(
				// TODO filter messages
				proc.SkipNoText,
			),
			proc.ForwardProcessor(),
		)
		return pipeline(state)
	})
}
