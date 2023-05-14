package xenv

import "os"

func Environ(env []string, filenames ...string) ([]string, error) { // TODO move this function somewhere else?
	for _, name := range filenames {
		cfgData, err := os.ReadFile(name)
		if err != nil {
			return nil, err
		}
		pairs, err := Parser([]rune(string(cfgData)))
		if err != nil {
			return nil, err
		}
		for _, v := range pairs {
			env = append(env, v[0]+"="+v[1]) // not 100% safe; side effects are possible
		}
	}
	return env, nil
}
