package xbot

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

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

func RequestMultipart(method string, to int64, fieldname string, data []byte, filename string) (*Request, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	err := w.WriteField("chat_id", strconv.FormatInt(to, 10))
	if err != nil {
		return nil, err
	}
	fw, err := w.CreateFormFile(fieldname, filename)
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
	Token  string
	Client Client
}

func (b *Bot) API(ctx context.Context, request *Request) ([]byte, error) {
	ctx = xlog.Ctx(ctx, "api", request.Method)
	err := error(nil)
	req := (*http.Request)(nil)
	resp := (*http.Response)(nil)
	data := []byte(nil)
	defer func() {
		xlog.Log(ctx, string(request.Body), string(data), err)
	}()
	reqUrl := "https://api.telegram.org/bot" + b.Token + "/" + request.Method
	req, err = http.NewRequestWithContext(ctx, http.MethodPost, reqUrl, bytes.NewReader(request.Body))
	if err != nil {
		return nil, xlog.Errorf(ctx, "request constructor: %w", err)
	}
	req.Header.Add("content-type", request.ContentType)
	resp, err = b.Client.Do(req)
	if err != nil {
		return nil, xlog.Errorf(ctx, "client: %w", err)
	}
	data, err = io.ReadAll(resp.Body)
	cerr := resp.Body.Close()
	if err != nil {
		if cerr != nil {
			return nil, xlog.Errorf(ctx, "reading: %w, and closing: %w", err, cerr)
		}
		return nil, xlog.Errorf(ctx, "reading: %w", err)
	}
	return data, nil
}
