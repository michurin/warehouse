package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func validateRoomID(rid string) error {
	if len(rid) > 30 {
		return fmt.Errorf("RoomID too long: %q", rid)
	}
	return nil
}

func validateID(id int64) error {
	if id < 0 {
		return fmt.Errorf("LastID has to be positive: %d", id)
	}
	return nil
}

// CustomValidator interface to implement you custom validation
// that performs before publish message
type CustomValidator interface {
	Validate(*http.Request, json.RawMessage) error
}

// ValidatorFunc is helper to turn function to CustomValidator interface
type ValidatorFunc func(*http.Request, json.RawMessage) error

func (f ValidatorFunc) Validate(r *http.Request, v json.RawMessage) error {
	return f(r, v)
}
