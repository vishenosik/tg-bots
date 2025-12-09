package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	botapi "github.com/vishenosik/tg-bots/internal/controller/bot"
	fileservice "github.com/vishenosik/tg-bots/internal/infra/file_service"
	"github.com/vishenosik/tg-bots/internal/store/csv"
	"github.com/vishenosik/tg-bots/internal/usecase"

	"github.com/vishenosik/gocherry"
	"github.com/vishenosik/gocherry/pkg/logs"
	bot "github.com/vishenosik/tg-bot-engine"

	_ctx "github.com/vishenosik/gocherry/pkg/context"
)

func main() {

	log := logs.SetupLogger().With(logs.AppComponent("main"))

	gocherry.Flags(os.Stdout, os.Args[1:],
		gocherry.AppFlags(os.Stdout),
		gocherry.ConfigFlags(os.Stdout),
	)

	flag.Parse()

	ctx := context.Background()

	app, err := NewApp(ctx)
	if err != nil {
		log.Error("failed to init app", logs.Error(err))
		os.Exit(1)
	}

	err = app.Start(ctx)
	if err != nil {
		log.Error("failed to start app", logs.Error(err))
		os.Exit(1)
	}

	// Graceful shut down
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	stopctx, stopCancel := context.WithTimeout(_ctx.WithStopCtx(context.Background(), <-stop), time.Second*5)
	defer stopCancel()

	app.Stop(stopctx)
}

func NewApp(ctx context.Context) (*gocherry.App, error) {

	// STORES

	fileservice, err := fileservice.NewFileService()
	if err != nil {
		return nil, err
	}

	csvStore := csv.NewCSVStore("senders.csv")

	// USECASES

	photousecase := usecase.NewPhotoUsecase(fileservice, csvStore)

	// API

	photosaver := botapi.NewPhotoApi(photousecase)
	affirmations := botapi.NewAffirmationsApi()

	// SERVICES

	telegramBot, err := bot.NewBotAPIEnv()
	if err != nil {
		return nil, err
	}

	photosaver.Route(telegramBot)
	affirmations.Route(telegramBot)

	app, err := gocherry.NewApp()
	if err != nil {
		return nil, err
	}

	app.AddServices(
		telegramBot,
	)

	return app, nil
}
