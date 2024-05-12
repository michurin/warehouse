package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/michurin/cnbot/app"
	"github.com/michurin/cnbot/xlog"
)

var Build = "development"

func main() {
	vFlag := flag.Bool("v", false, "show version")
	flag.Parse()
	if vFlag != nil && *vFlag {
		app.ShowVersionInfo()
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	app.SetupLogging()
	cfg, err := app.LoadConfigs(flag.Args()...)
	if err != nil {
		xlog.L(ctx, err)
		return
	}
	err = app.Application(ctx, cfg, Build+" "+app.MainVersion())
	if err != nil {
		xlog.L(ctx, err)
		return
	}
}
