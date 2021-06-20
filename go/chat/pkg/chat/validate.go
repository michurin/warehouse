package chat

import "fmt"

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
