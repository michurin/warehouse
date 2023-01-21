package smartjson

import (
	"fmt"
	"strings"
)

// TODO move Str to separate package

type Str struct {
	str string
	l   int
}

func New(s string) Str {
	return Str{
		str: s,
		l:   len([]rune(s)),
	}
}

func (s Str) String() string {
	return s.str
}

func (s Str) Len() int {
	return s.l
}

func Concat(a ...Str) Str {
	strs := make([]string, len(a))
	l := 0
	for i, v := range a {
		strs[i] = v.str
		l += v.l
	}
	return Str{
		str: strings.Join(strs, ""),
		l:   l,
	}
}

func Join(a []Str, sep Str) Str {
	strs := make([]string, len(a)*2-1) // TODO it'll be great to check len overflow
	l := 0
	for i, v := range a {
		strs[i*2] = v.str
		l += v.l
		if i > 0 {
			strs[i*2-1] = sep.str
			l += sep.l
		}
	}
	return Str{
		str: strings.Join(strs, ""),
		l:   l,
	}
}

func Repeat(s string, n int) Str {
	return Str{
		str: strings.Repeat(s, n),
		l:   n * len([]rune(s)),
	}
}

// TODO move colors to separate package?
// TODO background?
// TODO bold, underline, italic...
// TODO 256-colors escapes
// TODO 16-colors escapes

var reset = Str{str: "\x1b[0m"}

func RGB(r, g, b int) Str {
	return Str{str: fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)}
}

func Wrap(s, color Str) Str {
	return Concat(color, s, reset)
}

// TODO move themes to separate package

// types: Map, Lst, Str, Flo, Fls, Tru, Num
// variant: Si, Mu
// element: Bo, Bc, Pa, Se
type Theme struct {
	MapSiBo Str // type — map
	MapSiBc Str // variant — (si)ngle-/(mu)lty- line
	MapSiPa Str
	MapSiSe Str
	MapMuBo Str
	MapMuBc Str
	MapMuPa Str
	MapMuSe Str

	LstSiBo Str
	LstSiBc Str
	LstSiSe Str
	LstMuBo Str
	LstMuBc Str
	LstMuSe Str

	Flo Str // for all scalars Mu is actual for overflow only
	Nul Str
	Fal Str
	Tru Str

	StrKeyQuo  Str
	StrKeyBody Str
	StrValQuo  Str
	StrValBody Str

	ErrorQuo  Str
	ErrorBody Str
}

var ThemeOne = Theme{
	MapSiBo: Wrap(New("{"), RGB(0, 255, 0)),
	MapSiBc: Wrap(New("}"), RGB(0, 255, 0)),
	MapSiPa: Wrap(New(":"), RGB(0, 255, 0)),
	MapSiSe: Wrap(New(","), RGB(0, 255, 0)),
	MapMuBo: Wrap(New("{"), RGB(255, 0, 0)),
	MapMuBc: Wrap(New("}"), RGB(255, 0, 0)),
	MapMuPa: Wrap(New(": "), RGB(255, 0, 0)),
	MapMuSe: Wrap(New(","), RGB(255, 0, 0)),

	LstSiBo: Wrap(New("["), RGB(0, 255, 0)),
	LstSiBc: Wrap(New("]"), RGB(0, 255, 0)),
	LstSiSe: Wrap(New(","), RGB(0, 127, 0)),
	LstMuBo: Wrap(New("["), RGB(255, 0, 0)),
	LstMuBc: Wrap(New("]"), RGB(255, 0, 0)),
	LstMuSe: Wrap(New(","), RGB(127, 0, 0)),

	Flo: RGB(0, 255, 255),
	Nul: Wrap(New("null"), RGB(127, 0, 0)),
	Fal: Wrap(New("false"), RGB(255, 0, 0)),
	Tru: Wrap(New("true"), RGB(0, 255, 0)),

	StrKeyQuo:  Wrap(New(`"`), RGB(0, 255, 255)),
	StrKeyBody: RGB(0, 127, 127),
	StrValQuo:  Wrap(New(`"`), RGB(255, 255, 0)),
	StrValBody: RGB(127, 127, 0),

	ErrorQuo:  Wrap(New(`"`), RGB(255, 0, 0)),
	ErrorBody: RGB(127, 0, 0),
}
