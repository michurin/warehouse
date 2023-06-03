package xctrl

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	xlog "github.com/michurin/minlog"

	"github.com/michurin/cnbot/app"
	"github.com/michurin/cnbot/xbot"
	"github.com/michurin/cnbot/xproc"
)

func Handler(bot *xbot.Bot, cmd *xproc.Cmd, loggingPatch xlog.LogPatch) http.HandlerFunc { //nolint:gocognit // reason to refactor
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := xlog.ApplyPatch(r.Context(), loggingPatch)
		// TODO mark ctx for logging?
		// TODO put http method to ctx
		// TODO put http content-type to ctx
		body, err := io.ReadAll(r.Body)
		if err != nil {
			app.Log(ctx, "body reading:", err)
		}
		method := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
		data := []byte(nil)
		switch r.Method {
		case http.MethodGet:
			data, err = bot.API(ctx, &xbot.Request{Method: method})
		case http.MethodPost:
			ct := r.Header.Get("content-type")
			if ct == "application/json" || strings.Contains(ct, "multipart/form-data") {
				data, err = bot.API(ctx, &xbot.Request{
					Method:      method,
					ContentType: ct,
					Body:        body,
				})
			} else {
				var to int64          // TODO refactor
				var req *xbot.Request // TODO refactor
				to, err = strconv.ParseInt(r.URL.Query().Get("to"), 10, 64)
				if err != nil {
					app.Log(ctx, err) // TODO response!
					return
				}
				// TODO add `to` to log context
				req, err = xbot.RequestFromBinary(body, to)
				if err != nil {
					app.Log(ctx, err) // TODO response!
					return
				}
				data, err = bot.API(ctx, req)
			}
		case "RUN":
			q := r.URL.Query()
			to, err := strconv.ParseInt(r.URL.Query().Get("to"), 10, 64) // TODO code dup
			if err != nil {
				app.Log(ctx, err) // TODO response!
				return
			}
			ctx := xlog.Ctx(ctx, "user", to)
			lPatch := xlog.Patch(ctx)
			go func() { // TODO: limit concurrency
				ctx := xlog.ApplyPatch(context.Background(), lPatch)
				// TODO refactor. it is similar to processMessage
				body, err := cmd.Run(ctx, q["a"], nil)
				if err != nil {
					app.Log(ctx, err)
					return
				}
				req, err := xbot.RequestFromBinary(body, to)
				if err != nil {
					app.Log(ctx, err)
					return
				}
				if req == nil { // TODO hmm... it happens?
					app.Log(ctx, "Script response skipped")
					return
				}
				_, err = bot.API(ctx, req) // TODO check body?
				if err != nil {
					app.Log(ctx, err)
					return
				}
			}()
			return
		default:
			app.Log(ctx, fmt.Errorf("method not allowed: "+r.Method))
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if err != nil {
			app.Log(ctx, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		// TODO consider `silent=true` parameter and skip writing if present
		_, err = w.Write(data) // TODO consider error
		if err != nil {
			app.Log(ctx, err)
			return
		}
	}
}
