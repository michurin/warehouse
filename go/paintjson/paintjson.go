package paintjson

const (
	outOfString = iota
	inQStr
	inQStrEscaped
	inNotStr
)

var (
	clrKey = []rune("\033[1;33m")
	clrS   = []rune("\033[1;36m")
	clrCtl = []rune("\033[1;31m")
	clrOff = []rune("\033[0m")
)

func PJ(s string) string {
	out := []rune(nil)
	state := outOfString
	lastWord := []rune(nil)
	lastSpaces := []rune(nil)
	for _, c := range s {
		switch state {
		case outOfString:
			switch c {
			case '{', '}', '[', ']', ':', ',':
				if lastWord != nil {
					if c == ':' { // it is key
						out = append(out, clrKey...)
						out = append(out, lastWord...)
						out = append(out, clrOff...)
					} else { // it is ordinary string
						out = append(out, lastWord...)
					}
					lastWord = nil
				}
				if lastSpaces != nil {
					out = append(out, lastSpaces...)
					lastSpaces = nil
				}
				out = append(out, clrCtl...)
				out = append(out, c)
				out = append(out, clrOff...)
			case '\x20', '\n', '\r', '\t':
				if lastWord == nil {
					out = append(out, c)
				} else {
					lastSpaces = append(lastSpaces, c)
				}
			case '"':
				lastWord = append(lastWord, c)
				state = inQStr
			default:
				out = append(out, clrS...)
				out = append(out, c)
				state = inNotStr
			}
		case inQStr:
			switch c {
			case '\\':
				lastWord = append(lastWord, c)
				state = inQStrEscaped
			case '"':
				lastWord = append(lastWord, c)
				state = outOfString
			default:
				lastWord = append(lastWord, c)
			}
		case inQStrEscaped:
			lastWord = append(lastWord, c)
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
	if state == inNotStr {
		out = append(out, clrOff...)
	}
	out = append(out, lastWord...)
	out = append(out, lastSpaces...)
	return string(out)
}
