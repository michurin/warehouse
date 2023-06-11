package xloop

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/michurin/minlog"

	"github.com/michurin/cnbot/app/aw"
	"github.com/michurin/cnbot/xbot"
	"github.com/michurin/cnbot/xjson"
	"github.com/michurin/cnbot/xproc"
)

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
				aw.Log(ctx, "Skip message", err)
				continue
			}
			if req == nil {
				continue
			}
			_, err = bot.API(ctx, req)
			if err != nil {
				aw.Log(ctx, err) // TODO process error
			}
		}
		offset++
	}
}

func getUpdates(ctx context.Context, bot *xbot.Bot, offset int64) ([]any, error) {
	req, err := xbot.RequestStruct("getUpdates", map[string]any{"offset": offset, "timeout": 30})
	if err != nil {
		return nil, minlog.Errorf(ctx, "cannot build request")
	}
	bytes, err := bot.API(ctx, req)
	if err != nil {
		return nil, minlog.Errorf(ctx, "api: %w", err) // TODO all returns are too hard?
	}
	data := any(nil)
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, minlog.Errorf(ctx, "unmarshal: %w", err)
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
		return nil, minlog.Errorf(ctx, "no user id: %w", err)
	}
	ctx = minlog.Ctx(ctx, "user", userID)
	env, err := xjson.JSONToEnv(m)
	if err != nil {
		return nil, minlog.Errorf(ctx, "cannot create env: %w", err)
	}
	text, err := xjson.String(m, "message", "text")
	if err != nil {
		aw.Log(ctx, err) // return nil, err // TODO callback_query...
	}
	args := textToArgs(text)
	data, err := command.Run(ctx, args, env)
	if err != nil {
		return nil, err // already wrapped with context
	}
	req, err := xbot.RequestFromBinary(data, userID)
	if err != nil {
		return nil, minlog.Errorf(ctx, "invalid data: %w", err)
	}
	if req == nil { // TODO hmm... it happens?
		return nil, minlog.Errorf(ctx, "cannot prepare request (nil): %w", err)
	}
	return req, nil
}

func textToArgs(text string) []string {
	a := strings.Fields(strings.ToLower(text))
	b := make([]string, len(a))
	for i, v := range a {
		if len(v) > 0 && v[0] == '/' {
			v = v[1:]
		}
		b[i] = v
	}
	return b
}
