package paintjson

const (
	outOfString = iota
	inQStr
	inQStrEscaped
	inNotStr
)

type FSM struct {
	clrKey     []byte
	clrSpecStr []byte
	clrStr     []byte
	clrCtl     []byte
	clrOff     []byte
	state      int
	colored    bool
	lastWord   []byte
	lastSpace  []byte
}

func NewFSM(opts ...Option) *FSM {
	fsm := &FSM{
		clrKey:     Yellow,
		clrSpecStr: Cyan,
		clrStr:     nil,
		clrCtl:     Red,
		clrOff:     Off,
		state:      outOfString,
		colored:    false,
		lastWord:   nil,
		lastSpace:  nil,
	}
	for _, o := range opts {
		o(fsm)
	}
	return fsm
}

func (fsm *FSM) on(a []byte, c []byte) []byte {
	if len(c) == 0 {
		return a
	}
	fsm.colored = true
	return append(a, c...)
}

func (fsm *FSM) off(a []byte) []byte {
	if !fsm.colored {
		return a
	}
	fsm.colored = false
	return append(a, fsm.clrOff...)
}

func (fsm *FSM) Next(c byte) []byte {
	out := []byte(nil)
	switch fsm.state {
	case outOfString:
		switch c {
		case '{', '}', '[', ']', ':', ',':
			if fsm.lastWord != nil {
				if c == ':' { // it is key
					out = fsm.on(out, fsm.clrKey)
					out = append(out, fsm.lastWord...)
					out = fsm.off(out)
				} else { // it is ordinary string
					out = fsm.on(out, fsm.clrStr)
					out = append(out, fsm.lastWord...)
					out = fsm.off(out)
				}
				fsm.lastWord = nil
			}
			if fsm.lastSpace != nil {
				out = append(out, fsm.lastSpace...)
				fsm.lastSpace = nil
			}
			out = fsm.on(out, fsm.clrCtl)
			out = append(out, c)
			out = fsm.off(out)
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
			out = fsm.on(out, fsm.clrSpecStr)
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
			out = fsm.off(out)
			out = fsm.on(out, fsm.clrCtl)
			out = append(out, c)
			out = fsm.off(out)
			fsm.state = outOfString
		case '\x20', '\n', '\r', '\t':
			out = fsm.off(out)
			out = append(out, c)
			fsm.state = outOfString
		default:
			out = append(out, c)
		}
	}
	return out
}

func (fsm *FSM) Finish() []byte {
	out := []byte(nil)
	out = fsm.off(out)
	out = append(out, fsm.lastWord...)
	out = append(out, fsm.lastSpace...)
	fsm.lastWord = nil
	fsm.lastSpace = nil
	fsm.state = outOfString
	return out
}
