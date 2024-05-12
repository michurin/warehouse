package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/michurin/cnbot/app"
	"github.com/michurin/cnbot/xlog"
)

var Build = "development"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	app.SetupLogging()
	cfg, err := app.LoadConfigs(os.Args[1:]...)
	if err != nil {
		xlog.L(ctx, err)
		return
	}
	err = app.Application(ctx, cfg, Build)
	if err != nil {
		xlog.L(ctx, err)
		return
	}
}
