package httpauthmw

type AuthChecker func(username, password string) bool
