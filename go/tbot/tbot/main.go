package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/michurin/warehouse/go/tbot/app"
	"github.com/michurin/warehouse/go/tbot/xbot"
	"github.com/michurin/warehouse/go/tbot/xlog"
	"github.com/michurin/warehouse/go/tbot/xproc"
)

func main() {
	// setup logging
	xlog.Fields = []string{"api", "pid", "user"}
	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		xlog.LabelInfo = "\033[32;1mINFO\033[0m"
		xlog.LabelError = "\033[31;1mERROR\033[0m"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	bot := &xbot.Bot{
		APIOrigin: "https://api.telegram.org",
		Token:     os.Getenv("BOT_TOKEN"), // TODO config!
		Client:    http.DefaultClient,
	}

	command := &xproc.Cmd{ // TODO rename xcmd or xproc
		InterruptDelay: time.Second,
		KillDelay:      time.Second,
		Command:        "./x.sh", // TODO config!
		Cwd:            ".",      // TODO config?
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return app.Loop(ctx, bot, command)
	})

	server := &http.Server{Addr: ":9999", Handler: app.Handler(bot)}
	eg.Go(func() error {
		<-ctx.Done()
		cx, stop := context.WithTimeout(context.Background(), time.Second)
		defer stop()
		return server.Shutdown(cx)
	})

	eg.Go(func() error {
		return server.ListenAndServe()
	})

	err := eg.Wait()
	fmt.Println("Exit reason:", err)
}
