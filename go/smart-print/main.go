package main

import (
	"context"
	"time"

	"mxxx/mxxx"
)

func aa() {
	mxxx.P(map[int]any{1: context.Background()})
}

func a() {
	mxxx.P("OK")
	aa()
}

func main() {
	a()
	mxxx.SLEEP(time.Minute, "ok", 1, context.Background())
}
