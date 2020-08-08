# Simple wrapper to http.RoundTripper to log all requests and responses

## Install

    go get github.com/michurin/warehouse/go/loghttpclient

## Examples

Simplest usage:

    client := &http.Client{Transport: &loghttpclient.RoundTripper{}}

After that you will see all requests and responses.

Logging functions and next RoundTripper are customizable.
More examples can be found in [documentation](https://pkg.go.dev/github.com/michurin/warehouse/go/loghttpclient) and in [tests](transport_test.go).

## Similar tools

- [ClientTrace](https://golang.org/pkg/net/http/httptrace/)
- [ReverseProxy](https://golang.org/pkg/net/http/httputil/#ReverseProxy)
