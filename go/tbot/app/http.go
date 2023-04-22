package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/michurin/warehouse/go/tbot/xbot"
	"github.com/michurin/warehouse/go/tbot/xjson"
	"github.com/michurin/warehouse/go/tbot/xlog"
)

type Options struct {
	Token string
}

func Handler(bot *xbot.Bot, messageChan chan<- any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			xlog.Log(ctx, "body reading:", err)
		}
		method := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
		data := []byte(nil)
		switch r.Method {
		case http.MethodGet:
			data, err = bot.API(ctx, &xbot.Request{Method: method})
		case http.MethodPost:
			data, err = bot.API(ctx, &xbot.Request{
				Method:      method,
				ContentType: r.Header.Get("content-type"),
				Body:        body,
			})
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if err != nil {
			xlog.Log(ctx, err)
			return
		}
		// TODO messageChan <- data (?)
		w.WriteHeader(http.StatusOK)
		w.Write(data) // TODO consider error
	}
}

func Loop(ctx context.Context, bot *xbot.Bot, messageChan chan<- any) error {
	offset := int64(0)
	for { // TODO extract each iteration into dedicated function; make it function/object wrapable and observable
		req, err := xbot.RequestStruct("getUpdates", map[string]any{"offset": offset, "timeout": 30})
		if err != nil {
			return xlog.Errorf(ctx, "cannot build request")
		}
		bytes, err := bot.API(ctx, req)
		if err != nil {
			return xlog.Errorf(ctx, "api: %w", err)
		}
		data := any(nil)
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			return xlog.Errorf(ctx, "unmarshal: %w", err)
		}
		ok, err := xjson.Bool(data, "ok")
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("response is not ok: %#v", data)
		}
		result, err := xjson.Slice(data, "result")
		if err != nil {
			return err
		}
		if len(result) == 0 { // we won't change offset if there is no new messages
			continue
		}
		offset = 0 // offset can be dropped
		for _, m := range result {
			u, err := xjson.Int(m, "update_id")
			if err != nil {
				return err
			}
			if u > offset {
				offset = u
			}
			messageChan <- m
		}
		offset++
	}
}
