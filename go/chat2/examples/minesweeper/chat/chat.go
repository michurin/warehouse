package chat

import (
	"encoding/json"
	"fmt"

	"github.com/michurin/warehouse/go/chat2/examples/minesweeper/valid"
	"github.com/michurin/warehouse/go/chat2/text"
)

// dirty oversimplification: we user the same DTO to
// - parse request
// - keep message in storage
// - send response
type dto struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Text  string `json:"text"`
}

func Validator(raw []byte) ([]byte, error) {
	in := dto{}
	err := json.Unmarshal(raw, &in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", string(raw), err)
	}
	if err := valid.Color(in.Color); err != nil {
		return nil, err
	}
	// TODO split check_text and sanitize_text routines
	return json.Marshal(&dto{
		Name:  text.SanitizeText(in.Name, 10, "[noname]"),      // TODO const!
		Text:  text.SanitizeText(in.Text, 1000, "[nomessage]"), // TODO const!
		Color: in.Color,
	})
}
