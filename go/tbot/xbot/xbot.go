package xbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/michurin/warehouse/go/tbot/xlog"
)

// --- TODO move Request

type Request struct {
	Method      string
	ContentType string
	Body        []byte
}

func RequestStruct(method string, x any) (*Request, error) {
	d, err := json.Marshal(x)
	if err != nil {
		return nil, err
	}
	return &Request{
		Method:      method,
		ContentType: "application/json",
		Body:        d,
	}, nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func RequestFromBinary(data []byte, userID int64) (*Request, error) {
	contentType := http.DetectContentType(data)
	switch {
	case strings.HasPrefix(contentType, "text/"):
		// TODO check length
		// TODO preformated: "entities":[{"offset":0,"length":len([]rune(string(data))),"type":"pre"}]
		return RequestStruct("sendMessage", map[string]any{"chat_id": userID, "text": string(data)})
	case strings.HasPrefix(contentType, "image/"): // TODO to limit image formats
		return reqMultipart("sendPhoto", userID, "photo", data, "image."+contentType[6:], contentType) // TODO naive way, it can be 'x-icon' for instance
	case strings.HasPrefix(contentType, "video/"): // TODO to limit video formats
		return reqMultipart("sendVideo", userID, "video", data, "video."+contentType[6:], contentType)
	case strings.HasPrefix(contentType, "audio/"): // it seems application/ogg is not fully supported; it requires OPUS encoding
		return reqMultipart("sendAudio", userID, "audio", data, "audio."+contentType[6:], contentType)
	default: // TODO hmm... application/* and font/*
		return reqMultipart("sendDocument", userID, "document", data, "document", contentType) // TODO extension!?
	}
}

func reqMultipart(method string, to int64, fieldname string, data []byte, filename string, ctype string) (*Request, error) { // TODO legacy
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	err := w.WriteField("chat_id", strconv.FormatInt(to, 10))
	if err != nil {
		return nil, err
	}
	fw, err := w.CreatePart(textproto.MIMEHeader{ // see implementation of CreateFormFile
		"Content-Disposition": {fmt.Sprintf(
			`form-data; name="%s"; filename="%s"`,
			quoteEscaper.Replace(fieldname),
			quoteEscaper.Replace(filename),
		)},
		"Content-Type": {ctype},
	})
	if err != nil {
		return nil, err
	}
	_, err = fw.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return &Request{
		Method:      method,
		ContentType: w.FormDataContentType(),
		Body:        body.Bytes(),
	}, nil
}

// --- /TODO

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

type Bot struct {
	APIOrigin string // injection to be testable
	Token     string
	Client    Client // injection to be observable // TODO move all logging into client middleware?
}

func (b *Bot) API(ctx context.Context, request *Request) ([]byte, error) {
	ctx = xlog.Ctx(ctx, "api", request.Method)
	err := error(nil)
	req := (*http.Request)(nil)
	resp := (*http.Response)(nil)
	data := []byte(nil)
	defer func() {
		xlog.Log(ctx, request.Body, data, err)
	}()
	reqURL := b.APIOrigin + "/bot" + b.Token + "/" + request.Method
	req, err = http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(request.Body))
	if err != nil {
		return nil, xlog.Errorf(ctx, "request constructor: %w", err)
	}
	req.Header.Set("content-type", request.ContentType)
	resp, err = b.Client.Do(req)
	if err != nil {
		return nil, xlog.Errorf(ctx, "client: %w", err)
	}
	defer resp.Body.Close()           // we are skipping error here
	data, err = io.ReadAll(resp.Body) // we have to read and close Body even for non-200 responses
	if err != nil {
		return nil, xlog.Errorf(ctx, "reading: %w", err)
	}
	return data, nil
}
