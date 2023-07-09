package httpauthmw

import "context"

type userNameKeyType int

const userNameKey userNameKeyType = iota

func UserName(ctx context.Context) string {
	s, _ := ctx.Value(userNameKey).(string)
	return s
}
