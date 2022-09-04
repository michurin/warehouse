package chat

import (
	"encoding/json"
	"fmt"

	"github.com/michurin/warehouse/go/chat2/examples/minesweeper/valid"
	"github.com/michurin/warehouse/go/chat2/text"
)

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
	return json.Marshal(&dto{
		Name:  text.SanitizeText(in.Name, 10, "[noname]"),
		Text:  text.SanitizeText(in.Text, 1000, "[nomessage]"),
		Color: in.Color,
	})
}
