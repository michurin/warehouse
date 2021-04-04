package paintjson

import "testing"

func TestPJ(t *testing.T) {
	clrQ = []rune("(Q]")
	clrS = []rune("(S]")
	clrCtl = []rune("(C]")
	clrOff = []rune("[O)")
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
		out:  `(C]{[O)(Q]"one"[O)(C]:[O)(S]1[O)(C]}[O)`,
	}, {
		name: "spaces",
		in:   ` { "one" : 12 } `,
		out:  ` (C]{[O) (Q]"one"[O) (C]:[O) (S]12[O) (C]}[O) `,
	}, {
		name: "escaped",
		in:   `"o\"ne"`,
		out:  `(Q]"o\"ne"[O)`,
	}, {
		name: "invalid_has_to_be_closed",
		in:   `"one`,
		out:  `(Q]"one[O)`,
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
