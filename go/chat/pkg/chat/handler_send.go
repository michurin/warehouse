package chat

import (
	"io/ioutil"
	"net/http"
)

type SendHandler struct {
	Storage *Storage
}

func (h *SendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Storage.Add(Message{string(body)}) // TODO validate message, TODO use JSON for text, nickname, color etc
	w.Write([]byte("OK"))
}
