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
	envs := sdenv.NewCollectsion()
	envs.PushStd(os.Environ())
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, minlog.Errorf(ctx, "reading: %w", err)
		}
		pairs, err := sdenv.Parser(data)
		if err != nil {
			return nil, minlog.Errorf(ctx, "parser: %w", err)
		}
		envs.Push(pairs)
	}
	return xcfg.Cfg(ctx, envs.CollectionStd()), nil
}
