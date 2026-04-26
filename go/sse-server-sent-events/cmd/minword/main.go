package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/michurin/minchat/internal/cleanup"
	"github.com/michurin/minchat/internal/handlerstat"
	"github.com/michurin/minchat/internal/middleware"
	"github.com/michurin/minchat/internal/router"
	"github.com/michurin/minchat/internal/xhouse"
	"github.com/michurin/minchat/internal/xlog"
)

func main() {
	const pollingTimeout = 599 * time.Second
	const inactiveTime = 10 * time.Second // TODO too short, for debugging only
	const chatAddr = ":7011"
	const statAddr = ":6060"

	xlog.Init()

	house := xhouse.New()

	appServer := &http.Server{Addr: chatAddr, Handler: middleware.Logging(router.Handler(house, pollingTimeout))}
	statServer := &http.Server{Addr: statAddr, Handler: handlerstat.New(house)}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	wg := new(sync.WaitGroup)

	go func() {
		<-stop
		cancel()
	}()

	go func() {
		cleanup.RevisionLoop(house, inactiveTime) // TODO graceful shutdown
		cancel()
	}()

	wg.Go(func() {
		err := appServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				slog.Info("App server stopped gracefully")
			} else {
				slog.Error("Listener error: " + err.Error())
			}
		}
		cancel()
	})

	wg.Go(func() {
		err := statServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				slog.Info("Stat server stopped gracefully")
			} else {
				slog.Error("Listener error: " + err.Error())
			}
		}
		cancel()
	})

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := appServer.Shutdown(ctx); err != nil {
		slog.Error("App server shutdown error: " + err.Error())
	}

	if err := statServer.Shutdown(ctx); err != nil {
		slog.Error("App server shutdown error: " + err.Error())
	}

	wg.Wait()
}
