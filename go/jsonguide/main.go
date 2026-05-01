package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"unicode/utf8"
)

const helpMessage = "" +
	"jsonguide [-c] [-s] [-h] <in.json >out.txt\n" +
	"  -c force colored output\n" +
	"  -s shallow output (do not dive into JSON-strings)\n" +
	"  -h help message"

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
	const contextSize = 20
	offset := r.dec.InputOffset()
	body := r.body.Bytes()
	bodyLen := int64(len(body))
	a := min(max(offset-contextSize, 0), bodyLen)
	for a > 0 && !utf8.RuneStart(body[a]) {
		a--
	}
	b := min(offset+contextSize, bodyLen)
	for b < bodyLen && !utf8.RuneStart(body[b]) {
		b++
	}
	return strings.ReplaceAll(string(body[a:b]), "\n", `\n`)
}

type colorTheme struct {
	errPre  string
	errPost string
	eqPre   string
	eqPost  string
	sepPre  string
	sepPost string
	keyPre  string
	keyPost string
}

var colored colorTheme

func init() {
	const off = "\033[0m"
	colored.errPre = "\033[91m"
	colored.errPost = off
	colored.eqPre = "\033[92m"
	colored.eqPost = off
	colored.sepPre = "\033[43;30m\033[2K"
	colored.sepPost = off
	colored.keyPre = "\033[93m"
	colored.keyPost = off
}

type writer struct {
	c   colorTheme
	out io.Writer
}

func (w *writer) msg(key, val string) {
	fmt.Fprintf(w.out, "%s%s%s %s=%s %s\n", w.c.keyPre, key, w.c.keyPost, w.c.eqPre, w.c.eqPost, val)
}

func (w *writer) err(scope, key, err string) {
	fmt.Fprintf(w.out, "%s%s: [%s] %s%s\n", w.c.errPre, key, scope, err, w.c.errPost)
}

func (w *writer) sep() {
	fmt.Fprintln(w.out, w.c.sepPre+"---"+w.c.sepPost)
}

func array(source *tokenReader, w *writer, deepLook bool, prefix string) bool {
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
		if value(source, w, deepLook, pfx) {
			return true
		}
		n++
	}
}

func object(source *tokenReader, w *writer, deepLook bool, prefix string) bool {
	sep := ""
	if prefix[len(prefix)-1] != '.' { // len(prefix) != 0, always
		sep = "."
	}
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
		if value(source, w, deepLook, prefix+sep+keyPath(key)) {
			return true
		}
		empty = false
	}
}

func value(source *tokenReader, w *writer, deepLook bool, prefix string) bool {
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
			return object(source, w, deepLook, prefix)
		case '[':
			return array(source, w, deepLook, prefix)
		}
	case string:
		if deepLook {
			if startNested([]byte(t), "|", w, deepLook, prefix) {
				return false
			}
			dst := make([]byte, base64.StdEncoding.DecodedLen(len(t)))
			n, err := base64.StdEncoding.Decode(dst, []byte(t))
			if err == nil && startNested(dst[:n], "#", w, deepLook, prefix) {
				return false
			}
		}
		w.msg(prefix, t)
		return false
	case bool, nil, float64:
		w.msg(prefix, fmt.Sprintf("%v (%T)", tkn, t))
		return false
	}
	w.err("value", prefix, fmt.Sprintf("Unknown token: %[1]v (%[1]T)", tkn))
	return true
}

func startNested(t []byte, sep string, w *writer, deepLook bool, prefix string) bool {
	if json.Valid(t) {
		d := &tokenReader{dec: json.NewDecoder(bytes.NewReader(t)), body: nil}
		tkn, err := d.token()
		if err == nil && (tkn == json.Delim('{') || tkn == json.Delim('[')) {
			d.unread(tkn)
			value(d, w, deepLook, prefix+" "+sep+" .")
			return true
		}
	}
	return false
}

func keyPath(key string) string {
	if isSimpleKey(key) {
		return key
	}
	return quotedKey(key)
}

func isSimpleKey(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !(c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9') {
			return false
		}
	}
	return true
}

var specialSymbols = map[rune][]rune{
	'\n': {'\\', 'n'},
	'\r': {'\\', 'r'},
	'\t': {'\\', 't'},
}

func quotedKey(s string) string {
	o := []rune(nil)
	for _, r := range s {
		if a, ok := specialSymbols[r]; ok {
			o = append(o, a...)
			continue
		}
		if r == '\\' || r == '"' {
			o = append(o, '\\')
		}
		o = append(o, r)
	}
	return `["` + string(o) + `"]`
}

func App(in io.Reader, out io.Writer, outStream fs.File, args []string) int {
	o, err := outStream.Stat()
	colorMode := err == nil && o.Mode()&os.ModeCharDevice > 0
	deepLook := true
	for _, o := range args {
		switch o {
		case "-h":
			fmt.Fprintln(out, helpMessage)
			return 0
		case "-c":
			colorMode = true
		case "-s":
			deepLook = false
		}
	}
	w := &writer{out: out}
	if colorMode {
		w.c = colored
	}
	body := new(bytes.Buffer)
	dec := json.NewDecoder(io.TeeReader(in, body))
	for {
		if value(&tokenReader{dec: dec, body: body}, w, deepLook, ".") {
			return 1
		}
		if dec.More() {
			w.sep()
		} else {
			return 0
		}
	}
}

func main() {
	os.Exit(App(os.Stdin, os.Stdout, os.Stdout, os.Args[1:]))
}
