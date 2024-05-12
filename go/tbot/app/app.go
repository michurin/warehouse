package app

import (
	"context"
	"net/http"
	"path"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/michurin/cnbot/ctxlog"
	"github.com/michurin/cnbot/xbot"
	"github.com/michurin/cnbot/xcfg"
	"github.com/michurin/cnbot/xctrl"
	"github.com/michurin/cnbot/xlog"
	"github.com/michurin/cnbot/xloop"
	"github.com/michurin/cnbot/xproc"
)

func bot(ctx context.Context, eg *errgroup.Group, cfg xcfg.Config, build string) {
	bot := &xbot.Bot{
		APIOrigin: "https://api.telegram.org",
		Token:     cfg.Token,
		Client:    http.DefaultClient,
	}

	envCommon := []string{"tg_x_ctrl_addr=" + cfg.ControlAddr, "tg_x_build=" + build}

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
		err := xloop.Loop(ctxlog.Add(ctx, "comp", "loop"), bot, command)
		if err != nil {
			return ctxlog.Errorfx(ctx, "polling loop: %w", err)
		}
		return nil
	})

	server := &http.Server{Addr: cfg.ControlAddr, Handler: xctrl.Handler(bot, commandLong, ctxlog.Patch(ctxlog.Add(ctx, "comp", "ctrl")))}
	eg.Go(func() error {
		<-ctx.Done()
		cx, stop := context.WithTimeout(context.Background(), time.Second)
		defer stop()
		return server.Shutdown(cx) //nolint:contextcheck
	})

	eg.Go(func() error {
		err := server.ListenAndServe()
		if err != nil {
			return ctxlog.Errorfx(ctx, "control server: %w", err)
		}
		return nil
	})
}

func Application(rootCtx context.Context, bots map[string]xcfg.Config, build string) error {
	if len(bots) == 0 {
		return ctxlog.Errorfx(rootCtx, "there is no configuration")
	}
	eg, ctx := errgroup.WithContext(rootCtx)
	for name, cfg := range bots {
		bot(ctxlog.Add(ctx, "bot", name), eg, cfg, build)
	}
	xlog.L(ctx, "Run. Build="+build)
	return eg.Wait()
}
