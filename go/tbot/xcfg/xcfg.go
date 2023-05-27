package xcfg

import (
	"context"
	"fmt"
	"os"
	"strings"

	xlog "github.com/michurin/minlog"

	"github.com/michurin/warehouse/go/tbot/app"
)

var (
	varSfxs = []string{
		"ctrl_addr",
		"token",
		"long_running_script", // order matters; to be greedy longer has to be at beginning
		"script",
	}
	varAllowedSfxs = strings.Join(varSfxs, ", ")
)

const (
	varPrefix    = "tb_"
	varPrefixLen = len(varPrefix)
)

type Config struct {
	ControlAddr       string
	Token             string
	Script            string
	LongRunningScript string
}

func Cfg(osEnviron []string) map[string]Config { //nolint:gocognit
	ctx := xlog.Ctx(context.Background(), "comp", "cfg")
	x := map[string]map[string]string{}
	for _, pair := range osEnviron {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			app.Log(ctx, fmt.Errorf("skipping %q: cannot find `=`", pair))
			continue
		}
		ek := strings.ToLower(kv[0])
		ev := kv[1]
		if len(ev) == 0 {
			app.Log(ctx, fmt.Errorf("skipping %q: value is empty", pair))
			continue
		}
		if !strings.HasPrefix(ek, varPrefix) {
			continue
		}
		sfxNotFound := true
		for _, sfx := range varSfxs {
			if strings.HasSuffix(ek, "_"+sfx) {
				sfxNotFound = false
				k := "default"
				if len(ek) > varPrefixLen+len(sfx)+1 {
					k = strings.ToLower(ek[varPrefixLen : len(ek)-1-len(sfx)])
				}
				t := x[k]
				if t == nil {
					t = map[string]string{}
				}
				if x, ok := t[sfx]; ok {
					app.Log(ctx, fmt.Errorf("overriding %q by %q: %q", x, ev, pair))
				}
				t[sfx] = ev
				x[k] = t
				break
			}
		}
		if sfxNotFound {
			app.Log(ctx, fmt.Errorf("skipping %q: has TB prefix, but wrong suffix. Allowed: %s", pair, varAllowedSfxs))
		}
	}
	res := map[string]Config{}
	for k, v := range x {
		if len(v) != 4 {
			app.Log(ctx, fmt.Errorf("skipping bot name %q: incomplete set of options", k))
			continue
		}
		c := Config{
			ControlAddr:       v[varSfxs[0]],
			Token:             v[varSfxs[1]],
			Script:            v[varSfxs[3]],
			LongRunningScript: v[varSfxs[2]],
		}
		if strings.HasPrefix(c.Token, "@") {
			x, err := os.ReadFile(c.Token[1:])
			if err != nil {
				app.Log(ctx, fmt.Errorf("skipping bot name %q: cannot get token from file: %q: %w", k, c.Token, err))
				continue
			}
			c.Token = strings.TrimSpace(string(x))
		}
		res[k] = c
	}
	return res
}
