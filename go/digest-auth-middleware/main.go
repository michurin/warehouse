package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

const timeFormat = "20060102150405"

func parseTokens(s string) (map[string]string, error) {
	tt := [][]byte(nil)
	state := 0
	for _, c := range []byte(s) {
		switch state {
		case 0: // skipping spaces
			switch c {
			case '\x20':
			case '"':
				tt = append(tt, []byte{})
				state = 2
			case ',', '=':
				tt = append(tt, []byte{c})
			default:
				tt = append(tt, []byte{c})
				state = 1
			}
		case 1: // not quoted string; TODO do we must to consider `"` and `\` as errors?
			switch c {
			case '\x20':
				state = 0
			case ',', '=':
				tt = append(tt, []byte{c})
				state = 0
			default:
				tt[len(tt)-1] = append(tt[len(tt)-1], c)
			}
		case 2: // quoted string
			switch c {
			case '\\':
				state = 3
			case '"':
				state = 0
			default:
				tt[len(tt)-1] = append(tt[len(tt)-1], c)
			}
		case 3:
			tt[len(tt)-1] = append(tt[len(tt)-1], c)
			state = 2
		}
	}
	if len(tt) < 1 {
		return nil, fmt.Errorf("empty header")
	}
	if string(tt[0]) != "Digest" {
		return nil, fmt.Errorf("invalid prefix: %q", tt[0])
	}
	tt = tt[1:] // remove prefix, keeps pairs only
	res := map[string]string{}
	for i, v := range tt {
		switch i % 4 {
		case 1:
			if len(v) != 1 || v[0] != '=' {
				return nil, fmt.Errorf("not `=` in token #%d: %q", i, string(v))
			}
		case 3:
			if len(v) != 1 || v[0] != ',' {
				return nil, fmt.Errorf("not `,` in token #%d: %q", i, string(v))
			}
		case 2:
			res[string(tt[i-2])] = string(v)
		}
	}
	return res, nil
}

func checkDigestHeader(h string) error {
	pairs, err := parseTokens(h)
	if err != nil {
		return err
	}

	{
		tm, err := time.Parse(timeFormat, pairs["opaque"])
		if err != nil {
			return err
		}
		log.Print("| ", tm, " | ", time.Now(), " | ", time.Since(tm))
		if time.Since(tm) > time.Minute { // TODO arg
			return fmt.Errorf("timeout")
		}
		// curl http://localhost:9999/ --digest -u u:123 -v
		password := "123" // password is function of username+realm
		method := "GET"
		if pairs["algorithm"] != "MD5" {
			return fmt.Errorf("invalid algorithm: %q", pairs["algorithm"])
		}
		resp := md5text(
			md5text(pairs["username"], pairs["realm"], password), // STORE
			pairs["nonce"],
			pairs["nc"],
			pairs["cnonce"],
			pairs["qop"],
			md5text(method, pairs["uri"]),
		)
		if resp != pairs["response"] {
			return fmt.Errorf("not matched")
		}
	}
	return nil
}

func md5text(x ...string) string {
	a := md5.Sum([]byte(strings.Join(x, ":")))
	return hex.EncodeToString(a[:])
}

func randString() string {
	b := make([]byte, 75)
	_, err := rand.Read(b)
	if err != nil {
		panic(err) // TODO
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func main() {
	tmpl, err := template.New("root").Parse("<body><pre>{{.Time}}\n<a href=/?{{.Rnd}}>{{.Rnd}}</a></pre></body>\n")
	if err != nil {
		panic(err)
	}
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		log.Print(auth)
		err := checkDigestHeader(auth)
		if err != nil {
			fmt.Println("===== ERROR:", err.Error())
			// https://datatracker.ietf.org/doc/html/rfc2617#section-3.2.1
			// https://datatracker.ietf.org/doc/html/rfc2617#section-4
			ah := fmt.Sprintf(`Digest realm="z", qop=auth, algorithm=MD5, nonce="%s", opaque="%s"`, randString(), time.Now().UTC().Format(timeFormat))
			w.Header().Set("WWW-Authenticate", ah)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = tmpl.Execute(w, struct {
			Time string
			Rnd  int64
		}{Time: time.Now().Format(time.DateTime), Rnd: time.Now().Unix()})
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			return
		}
	})
	http.ListenAndServe(":9999", h)
}
