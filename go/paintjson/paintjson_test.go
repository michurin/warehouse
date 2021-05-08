package paintjson

import "testing"

func TestString(t *testing.T) {
	cases := []struct {
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
		exp:  `(C]{[O)(Q]"one"[O)(C]:[O)(s]"two"[O)(C]}[O)`,
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
	}}
	optsColors := []Option{
		ClrKey([]byte("(Q]")),
		ClrStr([]byte("(s]")),
		ClrSpecStr([]byte("(S]")),
		ClrCtl([]byte("(C]")),
		ClrOff([]byte("[O)")),
	}
	optsNoColors := []Option{
		ClrKey(nil),
		ClrStr(nil),
		ClrSpecStr(nil),
		ClrCtl(nil),
		ClrOff(nil),
	}
	for _, c := range cases {
		c := c
		t.Run("color_"+c.name, func(t *testing.T) {
			o := String(c.in, optsColors...)
			if o != c.exp {
				t.Errorf("%s != %s", o, c.exp)
			}
		})
	}
	for _, c := range cases {
		c := c
		t.Run("no_color_"+c.name, func(t *testing.T) {
			o := String(c.in, optsNoColors...)
			if o != c.in {
				t.Errorf("%s != %s", o, c.in)
			}
		})
	}
}
