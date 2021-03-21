package paintjson

const (
	outOfString = iota
	inQStr
	inQStrEscaped
	inNotStr
)

var clrQ = []rune("\033[32m")
var clrS = []rune("\033[36m")
var clrCtl = []rune("\033[35m")
var clrOff = []rune("\033[0m")

func PJ(s string) string {
	out := []rune(nil)
	state := outOfString
	for _, c := range s {
		switch state {
		case outOfString:
			switch c {
			case '{', '}', '[', ']', ':', ',':
				out = append(out, clrCtl...)
				out = append(out, c)
				out = append(out, clrOff...)
			case '\x20', '\n', '\r', '\t':
				out = append(out, c)
			case '"':
				out = append(out, clrQ...)
				out = append(out, c)
				state = inQStr
			default:
				out = append(out, clrS...)
				out = append(out, c)
				state = inNotStr
			}
		case inQStr:
			switch c {
			case '\\':
				out = append(out, c)
				state = inQStrEscaped
			case '"':
				out = append(out, c)
				out = append(out, clrOff...)
				state = outOfString
			default:
				out = append(out, c)
			}
		case inQStrEscaped:
			out = append(out, c)
			state = inQStr
		case inNotStr:
			switch c {
			case '{', '}', '[', ']', ':', ',':
				out = append(out, clrOff...)
				out = append(out, clrCtl...)
				out = append(out, c)
				out = append(out, clrOff...)
				state = outOfString
			case '\x20', '\n', '\r', '\t':
				out = append(out, clrOff...)
				out = append(out, c)
				state = outOfString
			default:
				out = append(out, c)
			}
		}
	}
	if state != outOfString {
		out = append(out, clrOff...)
	}
	return string(out)
}
