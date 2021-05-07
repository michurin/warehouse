package paintjson

func String(s string, opts ...Option) string {
	out := []byte(nil)
	fsm := NewFSM(opts...)
	for _, c := range []byte(s) {
		out = append(out, fsm.Next(c)...)
	}
	out = append(out, fsm.Finish()...)
	return string(out)
}
