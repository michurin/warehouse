package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	// curl -qsk https://localhost:8443/
	err := http.ListenAndServeTLS(":8443", "server.crt", "server.key", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, strings.NewReader("ok\n"))
	}))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
