package paintjson

import "testing"

func TestPJ(t *testing.T) {
	clrQ = []rune{'A'}
	clrS = []rune{'B'}
	clrCtl = []rune{'C'}
	clrOff = []rune{'O'}
	for _, c := range []struct {
		name string
		in   string
		out  string
	}{{
		name: "empty",
		in:   "",
		out:  "",
	}, {
		name: "simple",
		in:   `{"one":1}`,
		out:  `C{OA"one"OC:OB1OC}O`,
	}, {
		name: "spaces",
		in:   ` { "one" : 12 } `,
		out:  ` C{O A"one"O C:O B12O C}O `,
	}, {
		name: "escaped",
		in:   `"o\"ne"`,
		out:  `A"o\"ne"O`,
	}, {
		name: "invalid_has_to_be_closed",
		in:   `"one`,
		out:  `A"oneO`,
	}} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			o := PJ(c.in)
			if o != c.out {
				t.Errorf("%q != %q", o, c.out)
			}
		})
	}
}
