package paintjson

import "testing"

func TestPJ(t *testing.T) {
	clrKey = []byte("(Q]")
	clrS = []byte("(S]")
	clrCtl = []byte("(C]")
	clrOff = []byte("[O)")
	for _, c := range []struct {
		name string
		in   string
		exp  string
	}{{
		name: "empty",
		in:   "",
		exp:  "",
	}, {
		name: "simple_value_token",
		in:   `{"one":true}`,
		exp:  `(C]{[O)(Q]"one"[O)(C]:[O)(S]true[O)(C]}[O)`,
	}, {
		name: "simple_value_string",
		in:   `{"one":"two"}`,
		exp:  `(C]{[O)(Q]"one"[O)(C]:[O)"two"(C]}[O)`,
	}, {
		name: "spaces",
		in:   ` { "one" : 12 } `,
		exp:  ` (C]{[O) (Q]"one"[O) (C]:[O) (S]12[O) (C]}[O) `,
	}, {
		name: "escaped",
		in:   `{"o\"ne":1}`,
		exp:  `(C]{[O)(Q]"o\"ne"[O)(C]:[O)(S]1[O)(C]}[O)`,
	}, {
		name: "invalid_has_to_be_closed",
		in:   `[1`,
		exp:  `(C][[O)(S]1[O)`,
	}} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			o := PJ(c.in)
			if o != c.exp {
				t.Errorf("%s != %s", o, c.exp)
			}
		})
	}
}
