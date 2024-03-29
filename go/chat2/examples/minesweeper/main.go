package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/michurin/warehouse/go/chat2/examples/minesweeper/chat"
	"github.com/michurin/warehouse/go/chat2/examples/minesweeper/game"
	"github.com/michurin/warehouse/go/chat2/httppost"
	"github.com/michurin/warehouse/go/chat2/stream"
)

func bindAddr() string {
	if len(os.Args) == 2 {
		return os.Args[1]
	}
	return ":8080"
}

type subRequestDTO struct {
	Bounds []uint64 `json:"b"`
}

type subResponseDTO struct {
	Bounds []uint64          `json:"b"`
	Chat   []json.RawMessage `json:"chat,omitempty"`
	Game   []json.RawMessage `json:"game,omitempty"`
}

type gameRequestDTO struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	CID   string `json:"cid"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func castToRawMessage(x [][]byte) []json.RawMessage {
	r := make([]json.RawMessage, len(x))
	for i, v := range x {
		r[i] = json.RawMessage(v)
	}
	return r
}

func main() {
	const chatStreanCapacity = 100
	const gameStreanCapacity = 100
	const arenaWidth = 12  // TODO const
	const arenaHeight = 12 // TODO const

	logger := log.Default()
	addr := bindAddr()
	http.Handle("/", http.FileServer(http.Dir("examples/minesweeper/htdocs")))

	chatStream := stream.New(chatStreanCapacity)
	gameStream := stream.New(gameStreanCapacity)

	arena := game.New()
	resDto, err := arena.Setup(arenaWidth, arenaHeight)
	if err != nil {
		panic(err)
	}
	gameStream.Put(resDto)

	http.HandleFunc("/pub_chat", httppost.Handler(logger, func(ctx context.Context, requestBody []byte) ([]byte, error) {
		data, err := chat.Validator(requestBody)
		if err != nil {
			return nil, err
		}
		chatStream.Put(data)
		return nil, nil
	}))

	http.HandleFunc("/pub_game", httppost.Handler(logger, func(ctx context.Context, requestBody []byte) ([]byte, error) {
		request := gameRequestDTO{}
		err := json.Unmarshal(requestBody, &request)
		if err != nil {
			return nil, fmt.Errorf("can not unmarshal game request: %w", err)
		}
		openData, err := arena.Open(request.X, request.Y, request.CID, request.Name, request.Color)
		if err != nil {
			if err == game.ErrorNoRoom {
				return []byte(`{"nr": true}`), nil
			}
			return nil, err
		}
		if openData != nil { // only if we have updates
			gameStream.Put(openData)
		}
		return nil, nil
	}))

	http.Handle("/sub", httppost.Handler(logger, func(ctx context.Context, requestBody []byte) ([]byte, error) {
		request := subRequestDTO{}
		err := json.Unmarshal(requestBody, &request)
		if err != nil {
			return nil, fmt.Errorf("sub: cannot unmarshal: %w", err)
		}
		var reqBoundChat, reqBoundGame uint64
		bounds := request.Bounds
		if len(bounds) == 2 {
			reqBoundChat = bounds[0]
			reqBoundGame = bounds[1]
		}
		select {
		case <-chatStream.Waiter(reqBoundChat):
			streamData, boundChat := chatStream.Updates(reqBoundChat)
			bodyRes, err := json.Marshal(subResponseDTO{
				Bounds: []uint64{boundChat, reqBoundGame},
				Chat:   castToRawMessage(streamData),
			})
			if err != nil {
				return nil, fmt.Errorf("sub: chat: cannot marshal: %w", err)
			}
			return bodyRes, nil
		case <-gameStream.Waiter(reqBoundGame):
			streamData, boundGame := gameStream.Updates(reqBoundGame)
			var gameResp []json.RawMessage
			if boundGame-reqBoundGame <= gameStreanCapacity { // negative is big positive
				gameResp = castToRawMessage(streamData)
			} else {
				dump, err := arena.Dump()
				if err != nil {
					return nil, err
				}
				gameResp = []json.RawMessage{dump}
			}
			bodyRes, err := json.Marshal(subResponseDTO{
				Bounds: []uint64{reqBoundChat, boundGame},
				Game:   gameResp,
			})
			if err != nil {
				return nil, fmt.Errorf("sub: game: cannot marshal: %w", err)
			}
			return bodyRes, nil
		case <-time.After(30 * time.Second):
			// https://datatracker.ietf.org/doc/html/draft-loreto-http-bidirectional-07#section-5.5
			// Several experiments have shown success with timeouts as high as 120
			// seconds, but generally 30 seconds is a safer value.
			return []byte("{}"), nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}))

	log.Printf("Listing on %s", addr)
	err = http.ListenAndServe(addr, nil)
	log.Printf(err.Error())
}
