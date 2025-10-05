package main

import (
	"context"

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
}
