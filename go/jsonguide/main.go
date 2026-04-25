package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type tokenReader struct {
	dec       *json.Decoder
	lastToken *json.Token
	body      *bytes.Buffer
}

func (r *tokenReader) token() (json.Token, error) {
	if r.lastToken != nil {
		t := *r.lastToken
		r.lastToken = nil
		return t, nil
	}
	return r.dec.Token()
}

func (r *tokenReader) unread(t json.Token) {
	r.lastToken = &t
}

func (r *tokenReader) errContext() string {
	offset := r.dec.InputOffset()
	b := r.body.Bytes()
	return strings.ReplaceAll(string(b[max(offset-20, 0):min(offset+20, int64(len(b)))]), "\n", `\n`) // TODO rune unsafe
}

type writer struct {
	errPre  string
	errPost string
	out     io.Writer
}

func (w *writer) msg(key, val string) {
	fmt.Fprintf(w.out, "%s = %s\n", key, val)
}

func (w *writer) err(scope, key, err string) {
	fmt.Fprintf(w.out, "%s%s: [%s] %s%s\n", w.errPre, key, scope, err, w.errPost)
}

func (w *writer) sep() {
	fmt.Fprintln(w.out, "---")
}

func array(source *tokenReader, w *writer, prefix string) bool {
	n := 0
	for {
		pfx := fmt.Sprintf("%s[%d]", prefix, n)
		tkn, err := source.token()
		if err == io.EOF {
			w.err("array", pfx, "Unexpected EOF")
			return true
		}
		if err != nil {
			w.err("array", pfx, "Parse error: ("+source.errContext()+") "+err.Error())
			return true
		}
		if t, ok := tkn.(json.Delim); ok {
			switch t {
			case ']':
				if n == 0 {
					w.msg(prefix, "[]")
				}
				return false
			case '}':
				w.err("array", pfx, "Unexpected delimiter")
				return true
			} // pass [, {
		}
		source.unread(tkn)
		if value(source, w, pfx) {
			return true
		}
		n++
	}
}

func object(source *tokenReader, w *writer, prefix string) bool {
	empty := true
	for {
		tkn, err := source.token()
		if err == io.EOF {
			w.err("object", prefix, "Unexpected EOF")
			return true
		}
		if err != nil {
			w.err("object", prefix, "Parse error: ("+source.errContext()+") "+err.Error())
			return true
		}
		if t, ok := tkn.(json.Delim); ok {
			if t == '}' {
				if empty {
					w.msg(prefix, "{}")
				}
				return false
			}
			w.err("object", prefix, "Unexpected delimiter")
			return true
		}
		key, ok := tkn.(string)
		if !ok {
			w.err("object", prefix, fmt.Sprintf("Key is not string: %[1]v (%[1]T)", tkn))
			return false
		}
		if value(source, w, prefix+"."+key) {
			return true
		}
		empty = false
	}
}

func value(source *tokenReader, w *writer, prefix string) bool {
	tkn, err := source.token()
	if err == io.EOF {
		w.err("value", prefix, "Unexpected EOF")
		return true
	}
	if err != nil {
		w.err("value", prefix, "Parse error: ("+source.errContext()+") "+err.Error())
		return true
	}
	switch t := tkn.(type) {
	case json.Delim:
		switch t {
		case '{':
			return object(source, w, prefix)
		case '[':
			return array(source, w, prefix)
		}
	case string:
		w.msg(prefix, t)
		return false
	case bool, nil, float64:
		w.msg(prefix, fmt.Sprintf("%v (%T)", tkn, t))
		return false
	}
	w.err("value", prefix, fmt.Sprintf("Unknown token: %[1]v (%[1]T)", tkn))
	return true
}

func App(in io.Reader, out io.Writer, isTerm bool) int {
	w := &writer{out: out}
	if isTerm {
		w.errPre = "\033[91m"
		w.errPost = "\033[0m"
	}
	body := new(bytes.Buffer)
	dec := json.NewDecoder(io.TeeReader(in, body))
	for {
		if value(&tokenReader{dec: dec, body: body}, w, "") {
			return 1
		}
		if dec.More() {
			w.sep()
		} else {
			return 0
		}
	}
}

func isTerminal(h *os.File) bool {
	o, err := h.Stat()
	return err == nil && o.Mode()&os.ModeCharDevice > 0
}

func main() {
	os.Exit(App(os.Stdin, os.Stdout, isTerminal(os.Stdout)))
}
