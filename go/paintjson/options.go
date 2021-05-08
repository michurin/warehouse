package paintjson

type Option func(fsm *FSM)

var (
	Yellow    = []byte("\033[33;1m")
	Brown     = []byte("\033[33m")
	Red       = []byte("\033[31;1m")
	Darkred   = []byte("\033[31m")
	Pink      = []byte("\033[35;1m")
	Darkpink  = []byte("\033[35m")
	Blue      = []byte("\033[34;1m")
	Darkblue  = []byte("\033[34m")
	Green     = []byte("\033[32;1m")
	Darkgreen = []byte("\033[32m")
	Cyan      = []byte("\033[36;1m")
	Darkcyan  = []byte("\033[36m")
	White     = []byte("\033[37;1m")
	Black     = []byte("\033[30m")
	Lightgray = []byte("\033[37m")
	Darkgray  = []byte("\033[30;1m")
	None      = []byte(nil)
	Off       = []byte("\033[0m")
)

func ClrKey(clr []byte) Option {
	return func(fsm *FSM) {
		fsm.clrKey = clr
	}
}

func ClrSpecStr(clr []byte) Option {
	return func(fsm *FSM) {
		fsm.clrSpecStr = clr
	}
}

func ClrStr(clr []byte) Option {
	return func(fsm *FSM) {
		fsm.clrStr = clr
	}
}

func ClrCtl(clr []byte) Option {
	return func(fsm *FSM) {
		fsm.clrCtl = clr
	}
}

func ClrOff(clr []byte) Option {
	return func(fsm *FSM) {
		fsm.clrOff = clr
	}
}
