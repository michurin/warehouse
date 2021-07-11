package main

import (
	"fmt"

	"github.com/michurin/minlog"
)

func logLineFormatter(tm, level, label, caller, msg string) string {
	c := "\033[32;1m"
	if level != minlog.DefaultInfoLabel {
		c = "\033[31;1m"
	}
	return fmt.Sprintf("%s %s%s\033[0m %s \033[33m%s\033[0m %s", tm, c, level, label, caller, msg)
}

func setupLogger() *minlog.Logger {
	return minlog.New(
		minlog.WithLabelPlaceholder("-"),
		minlog.WithLineFormatter(logLineFormatter),
	)
}
