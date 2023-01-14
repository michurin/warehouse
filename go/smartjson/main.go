package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type Opts struct { // TODO: colors, separators; colors have to be considered on width calculation(!)
	Width  int
	Indent int
}

var reprStringReplacer = strings.NewReplacer(`"`, `\"`, `/`, `\/`, "\b", `\b`, "\f", `\f`, "\n", `\n`, "\r", `\r`, "\t", `\t`, `\`, `\\`)

func reprString(s string) string {
	return `"` + reprStringReplacer.Replace(s) + `"`
}

func max(a int, b int) int {
	if b > a {
		return b
	}
	return a
}

func marshal(level int, inp any, opts *Opts) (string, string) {
	ppfx := strings.Repeat(" ", opts.Indent*level)
	level++
	pfx := strings.Repeat(" ", opts.Indent*level)
	switch val := inp.(type) {
	case map[string]any:
		kk := make([]string, 0, len(val))
		for k := range val {
			kk = append(kk, k)
		}
		sort.Strings(kk)
		aa := make([]string, len(kk)) // aa — strings with single line representation for short single line relust
		bb := make([]string, len(kk)) // bb — multi line result with multi line parts
		cc := make([]string, len(kk)) // cc — multi line result with single line parts
		f := true                     // flag — is full set
		w := 0
		for i, k := range kk {
			s, m := marshal(level, val[k], opts)
			f = f && s != ""
			aa[i] = reprString(k) + ": " + s
			bb[i] = pfx + reprString(k) + ": " + m
			cc[i] = pfx + reprString(k) + ": " + s
			w = max(w, len(cc[i]))
		}
		t := bb
		s := ""
		if f {
			if w < opts.Width {
				t = cc
			}
			s = "{" + strings.Join(aa, ", ") + "}"
		}
		return s, "{\n" + strings.Join(t, ", ") + "\n" + ppfx + "}"
	case []any:
		aa := make([]string, len(val)) // the same naming as above
		bb := make([]string, len(val))
		cc := make([]string, len(val))
		f := true
		w := 0
		for i, v := range val {
			s, m := marshal(level, v, opts)
			f = f && s != ""
			aa[i] = s
			bb[i] = pfx + m
			cc[i] = pfx + s
			w = max(w, len(cc[i]))
		}
		t := bb
		s := ""
		if f {
			if w < opts.Width {
				t = cc
			}
			s = "[" + strings.Join(aa, ", ") + "]"
		}
		return s, "[\n" + strings.Join(t, ",\n") + "\n" + ppfx + "]"
	case string:
		o := reprString(val)
		return o, o
	case float64:
		o := fmt.Sprintf("%g", val)
		return o, o
	default:
		panic(fmt.Sprintf("Unknown type %T", inp))
	}
}

func Marshal(inp any, opts *Opts) string {
	s, m := marshal(0, inp, opts)
	if len(s) < opts.Width {
		return s
	}
	return m
}

func consoleSize() (int, int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	s := string(out)
	s = strings.TrimSpace(s)
	sArr := strings.Split(s, " ")

	heigth, err := strconv.Atoi(sArr[0])
	if err != nil {
		return 0, 0, err
	}

	width, err := strconv.Atoi(sArr[1])
	if err != nil {
		return 0, 0, err
	}
	return heigth, width, nil
}

// echo '{"1": [1, {"x":"xx"}, 4444]}' | go run main.go
func main() {
	buff, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	target := any(nil)
	err = json.Unmarshal(buff, &target)
	if err != nil {
		panic(err)
	}
	fmt.Println(target)
	fmt.Println("====")
	singleLine, multiLine := marshal(0, target, &Opts{
		Width:  30,
		Indent: 4,
	})
	fmt.Println(singleLine)
	fmt.Println("12345678901234567890")
	fmt.Println(multiLine)

	fmt.Println("====")
	r := Marshal(target, &Opts{
		Width:  30,
		Indent: 4,
	})
	fmt.Println(r)
	w, h, err := consoleSize()
	fmt.Println(w, h, err)
}
