package loghttptransport

import "net/http"

type LogFunc func(r *http.Request, rErr error, rData []byte, respErr error, respData []byte)

func (f LogFunc) Log(r *http.Request, rErr error, rData []byte, respErr error, respData []byte) {
	f(r, rErr, rData, respErr, respData)
}

type logger interface {
	Log(r *http.Request, rErr error, rData []byte, respErr error, respData []byte)
}
