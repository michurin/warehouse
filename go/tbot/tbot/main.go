package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/michurin/warehouse/go/tbot/app"
	"github.com/michurin/warehouse/go/tbot/xbot"
	"github.com/michurin/warehouse/go/tbot/xcfg"
	"github.com/michurin/warehouse/go/tbot/xenv"
	"github.com/michurin/warehouse/go/tbot/xlog"
	"github.com/michurin/warehouse/go/tbot/xproc"
)

func setupLogging() {
	xlog.Fields = []xlog.Field{
		xlog.StdFieldTime,
		xlog.StdFieldLevel,
		{
			Name: "bot",
			Proc: func(a any) string {
				return a.(string)
			},
		},
		{
			Name: "comp",
			Proc: func(a any) string {
				return a.(string)
			},
		},
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
}

func bot(ctx context.Context, eg *errgroup.Group, cfg xcfg.Config) {
	bot := &xbot.Bot{
		APIOrigin: "https://api.telegram.org",
		Token:     cfg.Token,
		Client:    http.DefaultClient,
	}

	command := &xproc.Cmd{
		InterruptDelay: 5 * time.Second,
		KillDelay:      5 * time.Second,
		Command:        cfg.Script,
		Cwd:            path.Dir(cfg.Script),
	}

	commandLong := &xproc.Cmd{
		InterruptDelay: 10 * time.Minute,
		KillDelay:      10 * time.Second,
		Command:        cfg.LongRunningScript,
		Cwd:            path.Dir(cfg.LongRunningScript),
	}

	eg.Go(func() error {
		return app.Loop(ctx, bot, command)
	})

	server := &http.Server{Addr: cfg.ControlAddr, Handler: app.Handler(bot, commandLong)}
	eg.Go(func() error {
		<-ctx.Done()
		cx, stop := context.WithTimeout(context.Background(), time.Second)
		defer stop()
		return server.Shutdown(cx) //nolint:contextcheck
	})

	eg.Go(func() error {
		return server.ListenAndServe()
	})
}

func application(rootCtx context.Context, bots map[string]xcfg.Config) error {
	if len(bots) == 0 {
		return xlog.Errorf(rootCtx, "there is no configuration")
	}
	eg, ctx := errgroup.WithContext(rootCtx)
	for name, cfg := range bots {
		bot(xlog.Ctx(ctx, "bot", name), eg, cfg)
	}
	return eg.Wait()
}

func loadConfigs(files ...string) (map[string]xcfg.Config, error) {
	env, err := xenv.Environ(os.Environ(), files...)
	if err != nil {
		return nil, err
	}
	return xcfg.Cfg(env), nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	setupLogging()
	cfg, err := loadConfigs("tbot.env") // TODO hardcoded
	if err != nil {
		xlog.Log(ctx, err)
		return
	}
	err = application(ctx, cfg)
	if err != nil {
		xlog.Log(ctx, err)
		return
	}
}
