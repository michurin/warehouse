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
	xlog.Fields = []xlog.Field{
		xlog.StdFieldTime,
		xlog.StdFieldLevel,
		{
			Name: "api",
			Proc: func(a any) string {
				return a.(string)
			},
		},
		{
			Name: "pid",
			Proc: func(a any) string {
				return fmt.Sprintf("%v", a)
			},
		},
		{
			Name: "user",
			Proc: func(a any) string {
				return fmt.Sprintf("%v", a)
			},
		},
		xlog.StdFieldCaller,
		xlog.StdFieldOCaller,
		xlog.StdFieldMessage,
	}
	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		xlog.Fields[0].Proc = func(any) string {
			return "\033[1;34m" + time.Now().Format("2006-01-02 15:04:05.000") + "\033[0m"
		}
		xlog.Fields[1].Proc = func(x any) string {
			if x.(int) == xlog.LevelError {
				return "\033[1;33;41m[error]\033[0m"
			}
			return "\033[1;32m[info]\033[0m"
		}
		xlog.Fields[2].Proc = func(x any) string {
			return "\033[1;35m" + x.(string) + "\033[0m"
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	bot := &xbot.Bot{
		APIOrigin: "https://api.telegram.org",
		Token:     os.Getenv("BOT_TOKEN"), // TODO config!
		Client:    http.DefaultClient,
	}

	command := &xproc.Cmd{
		InterruptDelay: time.Second, // TODO config?
		KillDelay:      time.Second, // TODO config??
		Command:        "./x.sh",    // TODO config!
		Cwd:            ".",         // TODO config?
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
		return server.Shutdown(cx) //nolint:contextcheck
	})

	eg.Go(func() error {
		return server.ListenAndServe()
	})

	err := eg.Wait()
	fmt.Println("Exit reason:", err)
}
