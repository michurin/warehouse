package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// curl -qsvk 'https://localhost:8080/'

func main() {
	err := http.ListenAndServeTLS(":8080", "ssl/cert.pem", "ssl/key.pem", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%v\n", time.Now())
	}))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
