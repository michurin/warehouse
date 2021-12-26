package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type handler struct {
	next string
	done chan struct{}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := do(w, r, h.next)
	if err != nil {
		fmt.Println(err)
		close(h.done)
	}
}

func do(w http.ResponseWriter, r *http.Request, next string) error {
	ctx := r.Context()
	var err error
	body := []byte(nil)
	if r.Method == http.MethodPost {
		body, err = io.ReadAll(r.Body)
		if err != nil {
			return err
		}
	}
	client := http.Client{}
	url := *r.URL
	url.Scheme = "http"
	url.Host = next
	req, err := http.NewRequestWithContext(ctx, r.Method, url.String(), bytes.NewBuffer(body))
	req.Header.Set("accept-encoding", "*/*") // default is "gzip"
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	rbody, err := io.ReadAll(resp.Body)
	fmt.Println(string(rbody))
	if err != nil {
		return err
	}
	w.WriteHeader(resp.StatusCode)
	_, err = w.Write(rbody)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	done := make(chan struct{})
	server := &http.Server{
		Addr: ":3008",
		Handler: handler{
			next: "localhost:3000",
			done: done,
		},
	}
	go func() {
		<-done
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
