package app

import (
	"context"
	"os"

	"github.com/michurin/minlog"
	"github.com/michurin/systemd-env-file/sdenv"

	"github.com/michurin/cnbot/xcfg"
)

func LoadConfigs(files ...string) (map[string]xcfg.Config, error) {
	ctx := minlog.Ctx(context.Background(), "comp", "cfg")
	env, err := sdenv.Environ(os.Environ(), files...)
	if err != nil {
		return nil, minlog.Errorf(ctx, "configuration loading: %w", err)
	}
	return xcfg.Cfg(ctx, env), nil
}
