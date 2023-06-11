package app

import (
	"os"

	"github.com/michurin/systemd-env-file/sdenv"

	"github.com/michurin/cnbot/xcfg"
)

func LoadConfigs(files ...string) (map[string]xcfg.Config, error) {
	env, err := sdenv.Environ(os.Environ(), files...)
	if err != nil {
		return nil, err
	}
	return xcfg.Cfg(env), nil
}
