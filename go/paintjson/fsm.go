package paintjson

const (
	outOfString = iota
	inQStr
	inQStrEscaped
	inNotStr
	finished
)

type FSM struct {
	clrKey     []byte
	clrSpecStr []byte
	clrCtl     []byte
	clrOff     []byte
	state      int
	lastWord   []byte
	lastSpace  []byte
}

func NewFSM(opts ...Option) *FSM {
	fsm := &FSM{
		// TODO: +clrStr
		clrKey:     Yellow,
		clrSpecStr: Cyan,
		clrCtl:     Red,
		clrOff:     Off,
		state:      outOfString,
		lastWord:   nil,
		lastSpace:  nil,
	}
	for _, o := range opts {
		o(fsm)
	}
	return fsm
}

func (fsm *FSM) Next(c byte) []byte {
	out := []byte(nil)
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
	case finished:
		panic("FSM.Next() is called after FSM.Tail()")
	}
	return out
}

func (fsm *FSM) Finish() []byte {
	out := []byte(nil)
	if fsm.state == inNotStr {
		out = append(out, fsm.clrOff...)
	}
	out = append(out, fsm.lastWord...)
	out = append(out, fsm.lastSpace...)
	fsm.state = finished
	return out
}
