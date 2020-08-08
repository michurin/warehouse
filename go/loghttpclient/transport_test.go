package loghttpclient_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/michurin/warehouse/go/loghttpclient"
)

func fatal(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func Example_roundTripper() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "Hello, client")
		fatal(err)
	}))
	defer ts.Close()

	client := http.Client{Transport: &loghttpclient.RoundTripper{
		Next: http.DefaultTransport,
		RequestFmt: func(request *http.Request, body string) {
			fmt.Printf("REQUEST: %s %q\n", request.Method, body)
		},
		ResponseFmt: func(response *http.Response, body string) {
			fmt.Printf("RESPONSE: %s %q\n", response.Status, body)
		},
	}}

	res, err := client.Get(ts.URL)
	fatal(err)
	greeting, err := ioutil.ReadAll(res.Body)
	fatal(err)
	err = res.Body.Close()
	fatal(err)
	fmt.Printf("BODY: %s\n", greeting)

	res, err = client.Post(ts.URL, "michurin/x-app", bytes.NewReader([]byte("data=ok")))
	fatal(err)
	greeting, err = ioutil.ReadAll(res.Body)
	fatal(err)
	err = res.Body.Close()
	fatal(err)
	fmt.Printf("BODY: %s\n", greeting)
	// Output:
	// REQUEST: GET "(nil body)"
	// RESPONSE: 200 OK "Hello, client\n"
	// BODY: Hello, client
	//
	// REQUEST: POST "data=ok"
	// RESPONSE: 200 OK "Hello, client\n"
	// BODY: Hello, client
}
