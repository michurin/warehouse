package xenv

// Env formats pairs as the os.Environ() does.
func Env(x [][2]string) []string {
	r := make([]string, len(x))
	for i, v := range x {
		r[i] = v[0] + "=" + v[1]
	}
	return r
}
