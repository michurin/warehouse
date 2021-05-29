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
	hdr := w.Header()
	hdr.Set("content-type", "application/json; charset=UTF-8")
	w.Write([]byte(`{}`)) // JSON
}
