package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/michurin/warehouse/go/chat2/httppost"
	"github.com/michurin/warehouse/go/chat2/stream"
	"github.com/michurin/warehouse/go/chat2/text"
)

// ----------------------------------------

type PointDTO struct {
	X int `json:"x"`
	Y int `json:"y"`
	V int `json:"v"`
}

type UserInfoDTO struct {
	Score int    `json:"score"`
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type OpenDTO struct {
	Type       string      `json:"t"`
	UsersTable UserInfoDTO `json:"u"`
	Points     []PointDTO  `json:"a"`
}

type UserInfo struct {
	score int
	id    int // Immutable. Starting from 1. 0 means no user. Everywhere
	name  string
	color string
}

// значения:
// - [0, 8] — count
// - 9 — mine
// - >10 — uid*10+v; as well as uid >= 1
type Arena struct {
	width  int
	height int
	closed int
	arena  [][]int
	users  map[string]*UserInfo
	mx     *sync.Mutex
}

func NewArena(w, h int) *Arena {
	a := [][]int(nil)
	closed := w * h
	for j := 0; j < h; j++ {
		t := []int(nil)
		for i := 0; i < w; i++ {
			v := 0
			if (i & j & 15) == 0 {
				v = 9
				closed--
			}
			t = append(t, v)
		}
		a = append(a, t)
	}
	for j := 0; j < h; j++ {
		j1 := j - 1
		j2 := j + 2
		if j1 < 0 {
			j1 = 0
		}
		if j2 > h {
			j2 = h
		}
		for i := 0; i < w; i++ {
			if a[j][i] == 9 {
				continue
			}
			i1 := i - 1
			i2 := i + 2
			if i1 < 0 {
				i1 = 0
			}
			if i2 > w {
				i2 = w
			}
			s := 0
			for q := j1; q < j2; q++ {
				for p := i1; p < i2; p++ {
					if a[q][p] == 9 {
						s++
					}
				}
			}
			a[j][i] = s
		}
	}
	for _, p := range a {
		s := ""
		for _, q := range p {
			s += fmt.Sprintf("\x1b[%sm%2d\x1b[0m", map[int]string{
				0: "30;1",
				1: "31",
				2: "32",
				3: "33",
				4: "34",
				5: "35",
				6: "36",
				7: "37",
				8: "32;1",
				9: "37;1",
			}[q], q)
		}
		fmt.Println(s)
	}
	return &Arena{
		width:  w,
		height: h,
		closed: closed,
		arena:  a,
		users:  map[string]*UserInfo{},
		mx:     new(sync.Mutex),
	}
}

// Open does one turn and returns marshaled
// - points
// - users table
// - and error
func (a *Arena) Open(x, y int, cid, name, color string) ([]byte, error) {
	a.mx.Lock()
	defer a.mx.Unlock()
	ui := a.users[cid]
	if ui == nil {
		n := len(a.users)
		if n >= 20 {
			return nil, nil // TODO room is fool
		}
		ui = &UserInfo{
			score: 0,
			id:    n + 1,
			name:  name,
			color: color,
		}
		a.users[cid] = ui
	} else {
		ui.name = name
		ui.color = color
	}
	if a.arena[y][x] >= 10 {
		return nil, nil
	}
	points := []PointDTO(nil)
	arenaDelta := ui.id * 10
	if a.arena[y][x] == 9 { // boom
		ui.score = 0
		a.arena[y][x] += arenaDelta
		points = []PointDTO{{
			X: x,
			Y: y,
			V: a.arena[y][x],
		}}
	} else { // regular opening
		stack := []int{x, y}
		for n := 2; n > 0; {
			x := stack[n-2]
			y := stack[n-1]
			n -= 2
			stack = stack[:n]
			if x >= 0 && x < a.width && y >= 0 && y < a.height && a.arena[y][x] <= 8 {
				if a.arena[y][x] == 0 {
					stack = append(stack, x-1, y-1, x-1, y, x-1, y+1, x, y-1, x, y+1, x+1, y-1, x+1, y, x+1, y+1)
					n += 16
				}
				a.arena[y][x] += arenaDelta
				a.closed--
				ui.score++
				points = append(points, PointDTO{
					X: x,
					Y: y,
					V: a.arena[y][x],
				})
			}
		}
	}
	// TODO if a.closed == 0
	openDto, err := json.Marshal(OpenDTO{
		Type: "o",
		UsersTable: UserInfoDTO{
			Score: ui.score,
			ID:    ui.id,
			Name:  ui.name,
			Color: ui.color,
		},
		Points: points,
	})
	if err != nil {
		return nil, err
	}
	return openDto, nil
}

// ----------------------------------------

var reColorStr = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

func validator(raw []byte) ([]byte, error) { // TODO slightly oversimplified approach; rewrite using DTOs
	in := map[string]string{}
	err := json.Unmarshal(raw, &in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", string(raw), err)
	}
	color := in["color"]
	if !reColorStr.MatchString(color) {
		return nil, errors.New("invalid color")
	}
	return json.Marshal(map[string]string{
		"name":  text.SanitizeText(in["name"], 10, "[noname]"),
		"text":  text.SanitizeText(in["text"], 1000, "[nomessage]"),
		"color": color,
	})
}

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
	const arenaWidth = 20
	const arenaHeight = 20

	logger := log.Default()
	addr := bindAddr()
	http.Handle("/", http.FileServer(http.Dir("examples/minesweeper/htdocs")))

	arena := NewArena(arenaWidth, arenaHeight)
	chatStream := stream.New(chatStreanCapacity)
	gameStream := stream.New(gameStreanCapacity)

	http.HandleFunc("/pub_chat", httppost.Handler(logger, func(ctx context.Context, requestBody []byte) ([]byte, error) {
		data, err := validator(requestBody)
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
		if err != nil {
			return nil, err
		}
		openData, err := arena.Open(request.X, request.Y, request.CID, request.Name, request.Color)
		if err != nil {
			return nil, err
		}
		gameStream.Put(openData)
		return nil, nil
	}))

	http.Handle("/sub", httppost.Handler(logger, func(ctx context.Context, requestBody []byte) ([]byte, error) {
		request := subRequestDTO{}
		err := json.Unmarshal(requestBody, &request)
		if err != nil {
			return nil, fmt.Errorf("sub: cannot unmarshal: %w", err)
		}
		var boundChat, boundGame uint64
		var streamData [][]byte
		bounds := request.Bounds
		if len(bounds) == 2 {
			boundChat = bounds[0]
			boundGame = bounds[1]
		}
		select {
		case <-chatStream.Waiter(boundChat):
			streamData, boundChat = chatStream.Updates(boundChat)
			bodyRes, err := json.Marshal(subResponseDTO{
				Bounds: []uint64{boundChat, boundGame},
				Chat:   castToRawMessage(streamData),
			})
			if err != nil {
				return nil, fmt.Errorf("sub: chat: cannot marshal: %w", err)
			}
			return bodyRes, nil
		case <-gameStream.Waiter(boundGame):
			streamData, boundGame = gameStream.Updates(boundGame)
			// TODO detect game resets
			// TODO detect stream gaps
			bodyRes, err := json.Marshal(subResponseDTO{
				Bounds: []uint64{boundChat, boundGame},
				Game:   castToRawMessage(streamData),
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
	err := http.ListenAndServe(addr, nil)
	log.Printf(err.Error())
}
