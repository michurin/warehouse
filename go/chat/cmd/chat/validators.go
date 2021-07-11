package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode"
)

func trivialValidator(v json.RawMessage) error {
	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		return err
	}
	return checkText("message", &s, 1000)
}

// {type: "chat-message", text: "one", nick: "Max Power", color: "red"}
// {type: "game-tune", x: 3, y: 2, fig: "o"}
// {type: "game-reset", size: 3}
type simplaPayload struct {
	Type  string  `json:"type"`
	Text  *string `json:"text"`
	Nick  *string `json:"nick"`
	Color *string `json:"color"`
	X     *int    `json:"x"`
	Y     *int    `json:"y"`
	Fig   *string `json:"fig"`
	Size  *int    `json:"size"`
}

func checkText(n string, s *string, l int) error {
	if s == nil {
		return fmt.Errorf("no %s", n)
	}
	if len(*s) == 0 {
		return fmt.Errorf("%s is empty", n)
	}
	if idx := strings.IndexFunc(*s, unicode.IsControl); idx >= 0 {
		return fmt.Errorf("control chars in %s: %q", n, *s)
	}
	if len(*s) > l {
		return fmt.Errorf("%s too long", n)
	}
	return nil
}

func checkXY(n string, c *int) error {
	if c == nil {
		return fmt.Errorf("no %s", n)
	}
	if *c < 0 || *c > 15 {
		return fmt.Errorf("%s out of range: %d", n, *c)
	}
	return nil
}

func simpleValidator(v json.RawMessage) error {
	var m simplaPayload
	if err := json.Unmarshal(v, &m); err != nil {
		return err
	}
	switch m.Type {
	case "chat-message":
		if err := checkText("text", m.Text, 1000); err != nil {
			return err
		}
		if err := checkText("nick", m.Nick, 50); err != nil {
			return err
		}
		if m.Color == nil {
			return errors.New("no color")
		}
		switch *m.Color {
		case "black", "red", "green":
		default:
			return fmt.Errorf("invalid color: %q", *m.Color)
		}
	case "game-tune":
		if err := checkXY("x", m.X); err != nil {
			return err
		}
		if err := checkXY("y", m.Y); err != nil {
			return err
		}
		if m.Fig == nil {
			return errors.New("no fig")
		}
		switch *m.Fig {
		case "x", "o":
		default:
			return fmt.Errorf("invalid figure: %q", *m.Fig)
		}
	case "game-reset":
		if m.Size == nil {
			return errors.New("no size")
		}
		switch *m.Size {
		case 3, 15:
		default:
			return fmt.Errorf("invalid size %d", *m.Size)
		}
	default:
		return fmt.Errorf("Invalid message type: %q", m.Type)
	}
	return nil
}
