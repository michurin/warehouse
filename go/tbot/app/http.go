package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/michurin/warehouse/go/tbot/xbot"
	"github.com/michurin/warehouse/go/tbot/xjson"
	"github.com/michurin/warehouse/go/tbot/xlog"
	"github.com/michurin/warehouse/go/tbot/xproc"
)

func Handler(bot *xbot.Bot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// TODO mark ctx for logging?
		// TODO put http method to ctx
		// TODO put http content-type to ctx
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
			ct := r.Header.Get("content-type")
			xlog.Log(ctx, "content-type:", ct) // TODO remove
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
					xlog.Log(ctx, err) // TODO response!
					return
				}
				// TODO add `to` to log context
				req, err = xbot.RequestFromBinary(body, to)
				if err != nil {
					xlog.Log(ctx, err) // TODO response!
					return
				}
				data, err = bot.API(ctx, req)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if err != nil {
			xlog.Log(ctx, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		// TODO consider `silent=true` parameter and skip writing if present
		w.Write(data) // TODO consider error
	}
}

func getUpdates(ctx context.Context, bot *xbot.Bot, offset int64) ([]any, error) {
	req, err := xbot.RequestStruct("getUpdates", map[string]any{"offset": offset, "timeout": 30})
	if err != nil {
		return nil, xlog.Errorf(ctx, "cannot build request")
	}
	bytes, err := bot.API(ctx, req)
	if err != nil {
		return nil, xlog.Errorf(ctx, "api: %w", err) // TODO all returns are too hard?
	}
	data := any(nil)
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, xlog.Errorf(ctx, "unmarshal: %w", err)
	}
	ok, err := xjson.Bool(data, "ok") // TODO xjson.True()?
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("response is not ok: %#v", data)
	}
	result, err := xjson.Slice(data, "result")
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Loop(ctx context.Context, bot *xbot.Bot, command *xproc.Cmd) error {
	offset := int64(0)
	for {
		result, err := getUpdates(ctx, bot, offset)
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
			// TODO refactor: get env, get args, run command
			req, err := processMessage(ctx, m, command)
			if err != nil {
				xlog.Log(ctx, err)
				continue
			}
			if req == nil {
				continue
			}
			_, err = bot.API(ctx, req)
			if err != nil {
				xlog.Log(ctx, err) // TODO process error
			}
		}
		offset++
	}
}

func message(m any) (any, error) { // TODO remove!
	for _, k := range []string{"message", "callback_query"} {
		val, ok, err := xjson.Any(m, k)
		if err != nil {
			return nil, err // TODO wrap, mention k in err message
		}
		if ok {
			return val, nil
		}
	}
	return nil, fmt.Errorf("payload for userID not found")
}

func userID(m any) (int64, error) { // TODO consider all types
	val, err := message(m)
	if err != nil {
		return 0, err
	}
	userID, err := xjson.Int(val, "from", "id")
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func processMessage(ctx context.Context, m any, command *xproc.Cmd) (*xbot.Request, error) {
	userID, err := userID(m)
	if err != nil {
		xlog.Log(ctx, "userID:", err)
		return nil, fmt.Errorf("no user id")
	}
	ctx = xlog.Ctx(ctx, "user", userID)
	env, err := xjson.JsonToEnv(m)
	if err != nil {
		xlog.Log(ctx, err)
		return nil, err // TODO wrap error
	}
	text, err := xjson.String(m, "message", "text") // TODO consider callback_query.message.text, callback_query.message.data?
	if err != nil {
		xlog.Log(ctx, err)
		// return nil, err // TODO callback_query...
	}
	args := strings.Fields(strings.ToLower(text))
	data, err := command.Run(ctx, args, env)
	if err != nil {
		xlog.Log(ctx, err)
		return nil, err
	}
	req, err := xbot.RequestFromBinary(data, userID)
	if err != nil {
		xlog.Log(ctx, err)
		return nil, err
	}
	if req == nil { // TODO hmm... it happens?
		xlog.Log(ctx, "Script response skipped")
		return nil, err
	}
	return req, nil
}
