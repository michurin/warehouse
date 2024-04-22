package xloop

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/michurin/cnbot/app/aw"
	"github.com/michurin/cnbot/ctxlog"
	"github.com/michurin/cnbot/xbot"
	"github.com/michurin/cnbot/xjson"
	"github.com/michurin/cnbot/xproc"
)

var consideringMessageTypes = []string{ // TODO it has to be tunable
	"callback_query",   // this strings are using in two places:
	"inline_query",     // in getUpdate and
	"message",          // in parsing function
	"message_reaction", // we assume we can get userID from any types using the same way. I'm not sure it works
	"poll",
	"poll_answer",
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
				aw.L(ctx, ctxlog.Errorf("skip message: %w", err))
				continue
			}
			if req == nil {
				continue
			}
			_, err = bot.API(ctx, req)
			if err != nil {
				aw.L(ctx, err) // TODO process error
			}
		}
		offset++
	}
}

func getUpdates(ctx context.Context, bot *xbot.Bot, offset int64) ([]any, error) {
	req, err := xbot.RequestStruct("getUpdates", map[string]any{
		"offset":          offset,
		"timeout":         30,
		"allowed_updates": consideringMessageTypes,
	})
	if err != nil {
		return nil, ctxlog.Errorfx(ctx, "cannot build request")
	}
	bytes, err := bot.API(ctx, req)
	if err != nil {
		return nil, ctxlog.Errorfx(ctx, "api: %w", err) // TODO all returns are too hard?
	}
	data := any(nil)
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, ctxlog.Errorfx(ctx, "unmarshal: %w", err)
	}
	ok, err := xjson.Bool(data, "ok") // TODO xjson.True()?
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ctxlog.Errorf("response is not ok: %#v", data)
	}
	result, err := xjson.Slice(data, "result")
	if err != nil {
		return nil, err
	}
	return result, nil
}

func userID(m any) (int64, error) { // TODO consider all types
	for _, bodyKey := range consideringMessageTypes { // TODO we are thinking all messages has the same structure related userID
		body, ok, err := xjson.Any(m, bodyKey)
		if err != nil {
			return 0, err // TODO wrap, mention k in err message
		}
		if ok {
			path := []string{"from", "id"}
			if bodyKey == "message_reaction" { // hakish
				path = []string{"user", "id"}
			}
			userID, err := xjson.Int(body, path...)
			if err != nil {
				return 0, ctxlog.Errorf("user not found: key=%s, path=%v: %w", bodyKey, path, err)
			}
			return userID, nil
		}
	}
	return 0, ctxlog.Errorf("body not found: consider %v", consideringMessageTypes)
}

func userText(m any) (string, error) { // TODO consider all types
	for _, bodyKey := range consideringMessageTypes { // TODO we are thinking all messages has the same structure related userID
		body, ok, err := xjson.Any(m, bodyKey)
		if err != nil {
			return "", err // TODO wrap, mention k in err message
		}
		if ok {
			if bodyKey == "message" { // hakish
				return xjson.String(body, "text")
			}
			if bodyKey == "callback_query" {
				return xjson.String(body, "data")
			}
			return bodyKey, nil
		}
	}
	return "", ctxlog.Errorf("body not found: consider %v", consideringMessageTypes)
}

func processMessage(ctx context.Context, m any, command *xproc.Cmd) (*xbot.Request, error) {
	userID, err := userID(m)
	if err != nil {
		return nil, ctxlog.Errorfx(ctx, "no user id: %w", err)
	}
	ctx = ctxlog.Add(ctx, "user", userID)
	env, err := xjson.JSONToEnv(m)
	if err != nil {
		return nil, ctxlog.Errorfx(ctx, "cannot create env: %w", err)
	}
	text, err := userText(m)
	if err != nil {
		aw.L(ctx, err) // return nil, err // TODO callback_query...
	}
	args := textToArgs(text)
	data, err := command.Run(ctx, args, env)
	if err != nil {
		return nil, err // already wrapped with context
	}
	req, err := xbot.RequestFromBinary(data, userID)
	if err != nil {
		return nil, ctxlog.Errorfx(ctx, "invalid data: %w", err)
	}
	if req == nil { // TODO hmm... it happens?
		return nil, ctxlog.Errorfx(ctx, "cannot prepare request (nil): %w", err)
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
