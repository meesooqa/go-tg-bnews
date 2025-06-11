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
	internaltg "github.com/meesooqa/go-tg-bnews/internal/telegram"
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
	flow := auth.NewFlow(internaltg.AuthFlow{}, auth.SendCodeOptions{})
	client := telegram.NewClient(appID, appHash, telegram.Options{
		DC:     2,
		DCList: dcs.Test(),
	})
	return client.Run(ctx, func(ctx context.Context) error {
		if err := client.Auth().IfNecessary(ctx, flow); err != nil {
			return err
		}
		api := client.API()
		tgs := internaltg.NewService(api, logger)

		channelFrom, err := tgs.GetChannel(ctx, "test_bbbolt_001")
		if err != nil {
			return err
		}
		channelTo, err := tgs.GetChannel(ctx, "test_bbbolt_002")
		if err != nil {
			return err
		}
		messages, err := tgs.GetMessages(ctx, channelFrom)
		if err != nil {
			return err
		}
		if len(messages) == 0 {
			return fmt.Errorf("no messages found in channel: %s", channelFrom.Username)
		}
		// TODO filter messages
		return tgs.ForwardMessages(ctx, messages, channelFrom, channelTo)
	})
}
