package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"

	"github.com/michurin/warehouse/go/chat2/examples/minesweeper/valid"
)

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

type ResetDTO struct {
	Width  int `json:"w"`
	Height int `json:"h"`
}

type OpenDTO struct {
	UsersTable []UserInfoDTO `json:"u,omitempty"`
	Points     []PointDTO    `json:"a,omitempty"`
	Field      [][]int       `json:"f,omitempty"`
	GameOver   bool          `json:"go,omitempty"`
	Reset      *ResetDTO     `json:"r,omitempty"`
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

func New() *Arena {
	return &Arena{
		mx: new(sync.Mutex),
	}
}

func (a *Arena) Setup(w, h int) ([]byte, error) { // TODO it seems, we do not need this args
	a.mx.Lock()
	defer a.mx.Unlock()
	return a.setup(w, h)
}

func (a *Arena) setup(w, h int) ([]byte, error) {
	ar := [][]int(nil)
	closed := w * h
	for j := 0; j < h; j++ {
		t := []int(nil)
		for i := 0; i < w; i++ {
			v := 0
			if /* (i & j & 15) == 0 */ /* p, q := i%7, j%7; p < 3 && q < 3 && !(p == 1 && q == 1) */ rand.Float32() < .27 {
				v = 9
				closed--
			}
			t = append(t, v)
		}
		ar = append(ar, t)
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
			if ar[j][i] == 9 {
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
					if ar[q][p] == 9 {
						s++
					}
				}
			}
			ar[j][i] = s
		}
	}
	for _, p := range ar {
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
	a.arena = ar
	a.width = w
	a.height = h
	a.closed = closed
	a.users = map[string]*UserInfo{}
	respDto, err := json.Marshal(OpenDTO{
		Reset: &ResetDTO{
			Width:  w,
			Height: h,
		},
	})
	if err != nil {
		return nil, err
	}
	return respDto, nil
}

// Open does one turn and returns marshaled
// - points
// - users table
// - and error
func (a *Arena) Open(x, y int, cid, name, color string) ([]byte, error) {
	a.mx.Lock()
	defer a.mx.Unlock()
	if err := valid.Open(a.width, a.height, x, y, cid, name, color); err != nil {
		return nil, err
	}
	if a.closed == 0 {
		return a.setup(a.width, a.height)
	}
	ui := a.users[cid]
	if ui == nil {
		n := len(a.users)
		if n >= 20 { // TODO const
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
	if a.arena[y][x] >= 10 { // already opened, no updates
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
			if x < 0 || x >= a.width || y < 0 || y >= a.height {
				continue
			}
			e := a.arena[y][x]
			if e > 8 {
				continue
			}
			if e == 0 {
				stack = append(stack, x-1, y-1, x-1, y, x-1, y+1, x, y-1, x, y+1, x+1, y-1, x+1, y, x+1, y+1)
				n += 16
			}
			a.closed--
			ui.score += e * e
			e += arenaDelta
			a.arena[y][x] = e
			points = append(points, PointDTO{
				X: x,
				Y: y,
				V: e,
			})
		}
	}
	respDto, err := json.Marshal(OpenDTO{
		UsersTable: []UserInfoDTO{{ // incremental update
			Score: ui.score,
			ID:    ui.id,
			Name:  ui.name,
			Color: ui.color,
		}},
		Points:   points,
		GameOver: a.closed == 0,
	})
	if err != nil {
		return nil, err
	}
	return respDto, nil
}

func (a *Arena) Dump() ([]byte, error) {
	a.mx.Lock() // TODO RLock?
	defer a.mx.Unlock()
	usersDto := make([]UserInfoDTO, 0, len(a.users))
	for _, ui := range a.users {
		usersDto = append(usersDto, UserInfoDTO{
			Score: ui.score,
			ID:    ui.id,
			Name:  ui.name,
			Color: ui.color,
		})
	}
	f := make([][]int, len(a.arena)) // filter out closed cells
	for i, p := range a.arena {
		g := make([]int, len(p))
		for j, q := range p {
			if q >= 10 {
				g[j] = q
			}
		}
		f[i] = g
	}
	respDto, err := json.Marshal(OpenDTO{
		UsersTable: usersDto,
		Field:      f,
		GameOver:   a.closed == 0,
	})
	if err != nil {
		return nil, err
	}
	return respDto, nil
}
