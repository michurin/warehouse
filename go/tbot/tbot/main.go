package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"time"

	xlog "github.com/michurin/minlog"
	xenv "github.com/michurin/systemd-env-file/sdenv"
	"golang.org/x/sync/errgroup"

	"github.com/michurin/cnbot/app"
	"github.com/michurin/cnbot/xbot"
	"github.com/michurin/cnbot/xcfg"
	"github.com/michurin/cnbot/xctrl"
	"github.com/michurin/cnbot/xloop"
	"github.com/michurin/cnbot/xproc"
)

var Build = "development"

func prefix(next xlog.FieldFunc, prefix string) xlog.FieldFunc { // TODO move it to xlog package?
	return func(r xlog.Record) string {
		t := next(r)
		if t == "" {
			return t
		}
		return prefix + next(r)
	}
}

func color(next xlog.FieldFunc, colorCode string) xlog.FieldFunc { // TODO move it to xlog package?
	return func(r xlog.Record) string {
		t := next(r)
		if t == "" {
			return t
		}
		return "\033[" + colorCode + "m" + next(r) + "\033[0m"
	}
}

func setupLogging() {
	_, file, _, _ := runtime.Caller(0)
	pfx := strings.TrimSuffix(file, "tbot/main.go")
	opts := []xlog.Option(nil)
	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		opts = []xlog.Option{
			xlog.WithFields(
				xlog.FieldLevel("", "\033[1;33;41m ERROR \033[0m"),
				color(xlog.FieldCaller(pfx), "1;34"),
				color(xlog.FieldErrorCaller(pfx), "1;31"),
				color(xlog.FieldNamed("comp"), "32"),
				color(xlog.FieldNamed("bot"), "35"),
				color(xlog.FieldNamed("api"), "1;35"),
				color(xlog.FieldNamed("user"), "1;32"),
				prefix(color(xlog.FieldNamed("pid"), "33"), "PID:"),
				xlog.FieldFallbackKV("api", "bot", "comp", "pid", "user"),
				xlog.FieldMessage()),
		}
	} else {
		opts = []xlog.Option{
			xlog.WithFields(
				xlog.FieldLevel("[I]", "[E]"),
				xlog.FieldCaller(pfx),
				xlog.FieldErrorCaller(pfx),
				xlog.FieldFallbackKV(),
				xlog.FieldMessage()),
		}
	}
	app.Log = xlog.New(opts...).Log
}

func bot(ctx context.Context, eg *errgroup.Group, cfg xcfg.Config) {
	bot := &xbot.Bot{
		APIOrigin: "https://api.telegram.org",
		Token:     cfg.Token,
		Client:    http.DefaultClient,
	}

	envCommon := []string{"tg_x_ctrl_addr=" + cfg.ControlAddr, "tg_x_build=" + Build}

	command := &xproc.Cmd{
		InterruptDelay: 10 * time.Second,
		KillDelay:      10 * time.Second,
		Env:            envCommon,
		Command:        cfg.Script,
		Cwd:            path.Dir(cfg.Script),
	}

	commandLong := &xproc.Cmd{
		InterruptDelay: 10 * time.Minute,
		KillDelay:      10 * time.Minute,
		Env:            envCommon,
		Command:        cfg.LongRunningScript,
		Cwd:            path.Dir(cfg.LongRunningScript),
	}

	eg.Go(func() error {
		return xloop.Loop(xlog.Ctx(ctx, "comp", "loop"), bot, command)
	})

	server := &http.Server{Addr: cfg.ControlAddr, Handler: xctrl.Handler(bot, commandLong, xlog.Patch(xlog.Ctx(ctx, "comp", "ctrl")))}
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
	app.Log(ctx, "Run. Build="+Build)
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
	cfg, err := loadConfigs(os.Args[1:]...)
	if err != nil {
		app.Log(ctx, err)
		return
	}
	err = application(ctx, cfg)
	if err != nil {
		app.Log(ctx, err)
		return
	}
}
