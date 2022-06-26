package text

import "unicode"

// SanitizeText truncates, removes ctrl chars and collapse spaces.
// If result is empty string, it returns third argument.
// It considers any non printable chars as spaces.
func SanitizeText(a string, n int, empty string) string {
	r := []rune(nil)
	s := false
	for _, c := range a {
		if len(r) >= n {
			break
		}
		if unicode.IsSpace(c) || !unicode.IsPrint(c) {
			s = len(r) > 0 // skip spaces at the begging
			continue
		}
		if s {
			if len(r) >= n-1 { // drop spaces at last position
				break
			}
			r = append(r, '\u0020')
			s = false
		}
		r = append(r, c)
	}
	if len(r) == 0 {
		return empty
	}
	return string(r)
}
