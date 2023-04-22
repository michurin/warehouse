package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/michurin/warehouse/go/tbot/app"
	"github.com/michurin/warehouse/go/tbot/xbot"
	"github.com/michurin/warehouse/go/tbot/xjson"
	"github.com/michurin/warehouse/go/tbot/xlog"
	"github.com/michurin/warehouse/go/tbot/xproc"
)

func message(m any) (any, error) {
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

func UserID(m any) (int64, error) {
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

func main() {
	// setup logging
	xlog.Fields = []string{"api", "pid", "user"}
	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		xlog.LabelInfo = "\033[32;1mINFO\033[0m"
		xlog.LabelError = "\033[31;1mERROR\033[0m"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	bot := &xbot.Bot{
		Token:  os.Getenv("BOT_TOKEN"), // TODO config!
		Client: http.DefaultClient,
	}

	command := xproc.Cmd{ // TODO rename xcmd or xproc
		InterruptDelay: time.Second,
		KillDelay:      time.Second,
		Command:        "./x.sh", // TODO config!
		Cwd:            ".",      // TODO config?
	}

	messageChan := make(chan any, 1000)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return app.Loop(ctx, bot, messageChan)
	})

	eg.Go(func() error {
		for {
			select {
			case m, ok := <-messageChan:
				if !ok {
					return fmt.Errorf("channel closed") // impossible, maybe panic would be even better
				}
				{ // to create private copy of ctx // TODO function
					userID, err := UserID(m)
					if err != nil {
						xlog.Log(ctx, "userID:", err)
						continue
					}
					ctx := xlog.Ctx(ctx, "user", userID)
					env, err := xjson.JsonToEnv(m)
					if err != nil {
						xlog.Log(ctx, err)
						return err // TODO wrap error
					}
					text, err := xjson.String(m, "message", "text") // TODO consider callback_query.message.text, callback_query.message.data?
					if err != nil {
						xlog.Log(ctx, err)
					}
					args := strings.Fields(strings.ToLower(text))
					data, err := command.Run(ctx, args, env)
					if err != nil {
						xlog.Log(ctx, err)
						continue
					}
					req := (*xbot.Request)(nil)
					contentType := http.DetectContentType(data)
					xlog.Log(ctx, contentType) // TODO remove
					switch {
					case strings.HasPrefix(contentType, "text/"):
						// TODO check length
						req, err = xbot.RequestStruct("sendMessage", map[string]any{"chat_id": userID, "text": string(data)})
					case strings.HasPrefix(contentType, "image/"):
						req, err = xbot.RequestMultipart("sendPhoto", userID, "photo", data, "image."+contentType[6:]) // TODO naive way, it can be 'x-icon' for instance
					case strings.HasPrefix(contentType, "video/"):
						req, err = xbot.RequestMultipart("sendVideo", userID, "video", data, "video."+contentType[6:])
					case strings.HasPrefix(contentType, "audio/"): // TODO +application/ogg? or consider ogg as voice?
						req, err = xbot.RequestMultipart("sendAudio", userID, "audio", data, "audio."+contentType[6:])
					default: // TODO hmm... application/* and font/*
						req, err = xbot.RequestMultipart("sendDocument", userID, "document", data, "document") // TODO extension!?
					}
					if err != nil {
						xlog.Log(ctx, err)
						continue
					}
					if req == nil {
						xlog.Log(ctx, "Script response skipped")
						continue
					}
					_, err = bot.API(ctx, req)
					if err != nil {
						xlog.Log(ctx, err) // TODO process error
					}
				}
			case <-ctx.Done():
				xlog.Log(ctx, ctx.Err())
				return ctx.Err()
			}
		}
	})

	server := &http.Server{Addr: ":9999", Handler: app.Handler(bot, messageChan)}
	eg.Go(func() error {
		<-ctx.Done()
		cx, stop := context.WithTimeout(context.Background(), time.Second)
		defer stop()
		return server.Shutdown(cx)
	})

	eg.Go(func() error {
		return server.ListenAndServe()
	})

	err := eg.Wait()
	fmt.Println("Exit reason:", err)
}
