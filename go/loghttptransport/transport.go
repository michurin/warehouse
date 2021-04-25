package loghttptransport

import (
	"net/http"

	"github.com/michurin/warehouse/go/readcloserwatcher"
)

type roundTripper struct {
	logger logger
}

func New(logger logger) *roundTripper {
	return &roundTripper{logger: logger}
}

func (rt *roundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	var rChan chan readcloserwatcher.Result
	var respChan chan readcloserwatcher.Result
	r.Body, rChan = readcloserwatcher.Watcher(r.Body)
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err == nil {
		resp.Body, respChan = readcloserwatcher.Watcher(resp.Body)
	}
	go func() { // Oh.. main() can be finished before it :-(
		var rR readcloserwatcher.Result
		var respR readcloserwatcher.Result
		for rChan != nil || respChan != nil {
			select {
			case rR = <-rChan:
				rChan = nil
			case respR = <-respChan:
				respChan = nil
			}
		}
		rt.logger.Log(r, rR.Err, rR.Octets, respR.Err, respR.Octets)
	}()
	return resp, err
}
