package main

import (
	"fmt"
	"net/http"

	"github.com/michurin/minlog"
)

func setupLogger() {
	minlog.SetDefaultLogger(minlog.New(
		minlog.WithLabelPlaceholder("-"),
		minlog.WithLineFormatter(func(tm, level, label, caller, msg string) string {
			c := "\033[32;1m"
			if level != minlog.DefaultInfoLabel {
				c = "\033[31;1m"
			}
			return fmt.Sprintf("%s %s%s\033[0m %s \033[33m%s\033[0m %s", tm, c, level, label, caller, msg)
		}),
	))
}

func NewWraper(label string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(minlog.Label(r.Context(), label)))
		})
	}
}
