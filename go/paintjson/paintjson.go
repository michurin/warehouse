package paintjson

const (
	outOfString = iota
	inQStr
	inQStrEscaped
	inNotStr
	closed
)

var (
	clrKey = []rune("\033[1;33m")
	clrS   = []rune("\033[1;36m")
	clrCtl = []rune("\033[1;31m")
	clrOff = []rune("\033[0m")
)

type FSM struct {
	clrKey     []rune
	clrSpecStr []rune
	clrCtl     []rune
	clrOff     []rune
	state      int
	lastWord   []rune
	lastSpace  []rune
}

func NewFSM() *FSM { // TODO: colors: NewFSM(opetions... Option)
	return &FSM{
		// TODO: use global colors for compatibility only, is to be rewritten
		clrKey:     clrKey, // []rune("\033[1;33m"),
		clrSpecStr: clrS,   // []rune("\033[1;36m"),
		clrCtl:     clrCtl, // []rune("\033[1;31m"),
		clrOff:     clrOff, // []rune("\033[0m"),
		state:      outOfString,
		lastWord:   nil,
		lastSpace:  nil,
	}
}

func (fsm *FSM) Next(c rune) []rune {
	out := []rune(nil)
	switch fsm.state {
	case outOfString:
		switch c {
		case '{', '}', '[', ']', ':', ',':
			if fsm.lastWord != nil {
				if c == ':' { // it is key
					out = append(out, fsm.clrKey...)
					out = append(out, fsm.lastWord...)
					out = append(out, fsm.clrOff...)
				} else { // it is ordinary string
					out = append(out, fsm.lastWord...)
				}
				fsm.lastWord = nil
			}
			if fsm.lastSpace != nil {
				out = append(out, fsm.lastSpace...)
				fsm.lastSpace = nil
			}
			out = append(out, fsm.clrCtl...)
			out = append(out, c)
			out = append(out, fsm.clrOff...)
		case '\x20', '\n', '\r', '\t':
			if fsm.lastWord == nil {
				out = append(out, c)
			} else {
				fsm.lastSpace = append(fsm.lastSpace, c)
			}
		case '"':
			fsm.lastWord = append(fsm.lastWord, c)
			fsm.state = inQStr
		default:
			out = append(out, fsm.clrSpecStr...)
			out = append(out, c)
			fsm.state = inNotStr
		}
	case inQStr:
		switch c {
		case '\\':
			fsm.lastWord = append(fsm.lastWord, c)
			fsm.state = inQStrEscaped
		case '"':
			fsm.lastWord = append(fsm.lastWord, c)
			fsm.state = outOfString
		default:
			fsm.lastWord = append(fsm.lastWord, c)
		}
	case inQStrEscaped:
		fsm.lastWord = append(fsm.lastWord, c)
		fsm.state = inQStr
	case inNotStr:
		switch c {
		case '{', '}', '[', ']', ':', ',':
			out = append(out, fsm.clrOff...)
			out = append(out, fsm.clrCtl...)
			out = append(out, c)
			out = append(out, fsm.clrOff...)
			fsm.state = outOfString
		case '\x20', '\n', '\r', '\t':
			out = append(out, fsm.clrOff...)
			out = append(out, c)
			fsm.state = outOfString
		default:
			out = append(out, c)
		}
	case closed:
		panic("FSM.Next() is called after FSM.Tail()")
	}
	return out
}

func (fsm *FSM) Tail() []rune {
	out := []rune(nil)
	if fsm.state == inNotStr {
		out = append(out, fsm.clrOff...)
	}
	out = append(out, fsm.lastWord...)
	out = append(out, fsm.lastSpace...)
	fsm.state = closed
	return out
}

// TODO move to file with helpers
func PJ(s string) string {
	out := []rune(nil)
	fsm := NewFSM()
	for _, c := range s {
		out = append(out, fsm.Next(c)...)
	}
	out = append(out, fsm.Tail()...)
	return string(out)
}
