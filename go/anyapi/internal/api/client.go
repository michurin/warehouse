package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func do(method, handler string, in []byte, debug bool) ([]byte, error) {
	body := io.Reader(http.NoBody)
	if in != nil {
		body = bytes.NewBuffer(in)
	}
	req, err := http.NewRequest(method, handler, body) //nolint:noctx // TODO: add context
	if err != nil {
		return nil, fmt.Errorf("http request building error: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println(string(respBody))
		return nil, fmt.Errorf("req: %s %q: %s; resp: %s; status: %s", method, handler, string(in), respBody, resp.Status) //nolint:err113
	}
	if debug {
		fmt.Printf("%s %q: %s -> %q: %s\n", req.Method, handler, string(in), resp.Status, respBody)
	}
	return respBody, nil
}

func PostStructRawResponse(handler string, in any, debug bool) ([]byte, error) {
	body, err := json.Marshal(in) // body, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return nil, err
	}
	respBody, err := do(http.MethodPost, handler, body, debug)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func PostStruct(handler string, in any, debug bool) (any, error) {
	body, err := PostStructRawResponse(handler, in, debug)
	if err != nil {
		return nil, err
	}
	result := any(nil)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetStruct(handler string, debug bool) (any, error) {
	body, err := do(http.MethodGet, handler, nil, debug)
	if err != nil {
		return nil, err
	}
	result := any(nil)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
