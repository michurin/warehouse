package main

import (
	"log"
	"net/http"
	"time"

	"github.com/michurin/minchat/internal/cleanup"
	"github.com/michurin/minchat/internal/middleware"
	"github.com/michurin/minchat/internal/router"
	"github.com/michurin/minchat/internal/xhouse"
	"github.com/michurin/minchat/internal/xlog"
)

func main() {
	const pollingTimeout = 599 * time.Second
	const inactiveTime = 10 * time.Second // TODO too short, for debugging only
	const chatAddr = ":7011"

	xlog.Init()

	house := xhouse.New()

	go cleanup.RevisionLoop(house, inactiveTime)

	err := http.ListenAndServe(chatAddr, middleware.Logging(router.Handler(house, pollingTimeout)))
	if err != nil {
		log.Printf("Listener error: %s", err.Error())
	}
}
