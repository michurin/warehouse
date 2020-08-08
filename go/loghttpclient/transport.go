package loghttpclient

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type RequestFmtFunc func(request *http.Request, body string)
type ResponseFmtFunc func(response *http.Response, body string)

type RoundTripper struct {
	Next        http.RoundTripper
	RequestFmt  RequestFmtFunc
	ResponseFmt ResponseFmtFunc
}

func defaultRequestFmt(r *http.Request, body string) {
	fmt.Printf("\u001B[1m>>>\u001B[0m \u001B[32;1m%s %s\u001B[0m %s\n", r.Method, r.URL.String(), body)
}

func defaultResponseFmt(r *http.Response, body string) {
	fmt.Printf("\u001B[1m<<<\u001B[0m \u001B[33;1m%s\u001B[0m %s\n", r.Status, body)
}

func fetchBody(s io.ReadCloser) (io.ReadCloser, string, error) {
	if s == nil {
		return nil, "(nil body)", nil
	}
	a, err := ioutil.ReadAll(s)
	if err != nil {
		return nil, "", err
	}
	err = s.Close()
	if err != nil {
		return nil, "", err
	}
	return ioutil.NopCloser(bytes.NewReader(a)), string(a), nil
}

func (t *RoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	var bodyString string
	next := http.DefaultTransport
	reqFmt := defaultRequestFmt
	respFmt := defaultResponseFmt
	if t.Next != nil {
		next = t.Next
	}
	if t.RequestFmt != nil {
		reqFmt = t.RequestFmt
	}
	if t.ResponseFmt != nil {
		respFmt = t.ResponseFmt
	}
	req.Body, bodyString, err = fetchBody(req.Body)
	if err != nil {
		return
	}
	reqFmt(req, bodyString)
	resp, err = next.RoundTrip(req)
	if err != nil {
		return
	}
	resp.Body, bodyString, err = fetchBody(resp.Body)
	if err != nil {
		return
	}
	respFmt(resp, bodyString)
	return
}
