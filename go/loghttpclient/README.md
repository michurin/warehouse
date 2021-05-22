# Simple wrapper to http.RoundTripper to log all requests and responses

## Install

    go get github.com/michurin/warehouse/go/loghttpclient

## Examples

Simplest usage:

    client := &http.Client{Transport: &loghttpclient.RoundTripper{}}

After that you will see all requests and responses.

Logging functions and next RoundTripper are customizable.
More examples can be found in [documentation](https://pkg.go.dev/github.com/michurin/warehouse/go/loghttpclient) and in [tests](transport_test.go).

## Deficiencies

It read whole body of requests, wraps it back to `ReadCloser` and
pass to next `RoundTripper`. Such heavy-handed intervention
obviously restricts abilities of client. It works well on
short requests only.

The second point is we log the data *that was sent* by server
not the data *the was receive* by client. It is possible
the client doesn't even started reading data. It is good
to debug server, but not for debugging your client.

However, you may be happy to see server's responses
*before* code of your client will start work.

Altogether, this solution is good for short and small HTTP
requests.

## Similar tools

- [ClientTrace](https://golang.org/pkg/net/http/httptrace/)
- [ReverseProxy](https://golang.org/pkg/net/http/httputil/#ReverseProxy)
