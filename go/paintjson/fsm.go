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
	balance    int
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
		balance:    0, // brackets balance
		lastWord:   nil,
		lastSpace:  nil,
	}
	for _, o := range opts {
		o(fsm)
	}
	return fsm
}

func (fsm *FSM) on(a []byte, c []byte) []byte {
	if len(c) == 0 || fsm.balance == 0 {
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

func (fsm *FSM) last(c []byte) []byte {
	out := []byte(nil)
	if fsm.lastWord != nil {
		out = fsm.on(out, c)
		out = append(out, fsm.lastWord...)
		out = fsm.off(out)
		fsm.lastWord = nil
	}
	if fsm.lastSpace != nil {
		out = append(out, fsm.lastSpace...)
		fsm.lastSpace = nil
	}
	return out
}

func (fsm *FSM) inc(c byte) {
	switch c {
	case '{', '[':
		fsm.balance++
	}
}

func (fsm *FSM) dec(c byte) {
	if fsm.balance <= 0 {
		return
	}
	switch c {
	case '}', ']':
		fsm.balance--
	}
}

func (fsm *FSM) Next(c byte) []byte {
	out := []byte(nil)
	switch fsm.state {
	case outOfString:
		switch c {
		case '{', '}', '[', ']', ',':
			fsm.inc(c)
			out = append(out, fsm.last(fsm.clrStr)...)
			out = fsm.on(out, fsm.clrCtl)
			out = append(out, c)
			out = fsm.off(out)
			fsm.dec(c)
		case ':': // we consider ':' as part of JSON only after string key
			color := fsm.lastWord != nil
			out = append(out, fsm.last(fsm.clrKey)...)
			if color {
				out = fsm.on(out, fsm.clrCtl)
			}
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
		case '{', '}', '[', ']', ',': // exclude ':' cause it can appear after string only
			fsm.inc(c)
			out = fsm.off(out)
			out = fsm.on(out, fsm.clrCtl)
			out = append(out, c)
			out = fsm.off(out)
			fsm.state = outOfString
			fsm.dec(c)
		case '\x20', '\n', '\r', '\t', ':': // ':' just interrupts unquoted string (invalid JSON)
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
	out = append(out, fsm.last(nil)...)
	fsm.state = outOfString
	fsm.balance = 0
	return out
}
