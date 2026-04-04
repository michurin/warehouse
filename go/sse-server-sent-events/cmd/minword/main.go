package main

import (
	"log"
	"net/http"

	"sse/handler"
	"sse/loggingmw"
	"sse/room"
)

func main() {
	house := room.New()
	// TODO start housekeeper
	err := http.ListenAndServe(":7011", loggingmw.MW(handler.Handler(house)))
	if err != nil {
		log.Printf("Listener error: %s", err.Error())
	}
}
