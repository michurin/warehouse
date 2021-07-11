package chat

import (
	"encoding/json"
	"fmt"
	"io"
)

type pollingRequest struct {
	RoomID string `json:"room"`
	ID     int64  `json:"id"`
}

type pollingResponse struct {
	Messages []json.RawMessage `json:"messages"`
	LastID   int64             `json:"lastID"`
}

type publishingRequest struct {
	RoomID  string          `json:"room"`
	Message json.RawMessage `json:"message"`
}

func DecodePollingRequest(body io.Reader) (string, int64, error) {
	req := new(pollingRequest)
	err := json.NewDecoder(body).Decode(req)
	if err != nil {
		return "", 0, err
	}
	if len(req.RoomID) > 30 {
		return "", 0, fmt.Errorf("RoomID too long: %q", req.RoomID)
	}
	if req.ID < 0 {
		return "", 0, fmt.Errorf("LastID has to be positive: %d", req.ID)
	}
	return req.RoomID, req.ID, nil
}

func DecodePublishingRequest(body io.Reader) (string, json.RawMessage, error) {
	req := new(publishingRequest)
	err := json.NewDecoder(body).Decode(req)
	if err != nil {
		return "", nil, err
	}
	if len(req.RoomID) > 30 {
		return "", nil, fmt.Errorf("RoomID too long: %q", req.RoomID)
	}
	return req.RoomID, req.Message, nil
}

func EncodePollingResponse(w io.Writer, mm []json.RawMessage, lastID int64) error {
	return json.NewEncoder(w).Encode(pollingResponse{
		Messages: mm,
		LastID:   lastID,
	})
}

func EncodePublishingResponse(w io.Writer) error {
	_, err := w.Write([]byte("{}"))
	return err
}
