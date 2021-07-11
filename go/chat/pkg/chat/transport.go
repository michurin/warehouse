package chat

import "net/http"

func SetContentType(hdr http.Header) {
	hdr.Set("content-type", "application/json; charset=UTF-8")
}
