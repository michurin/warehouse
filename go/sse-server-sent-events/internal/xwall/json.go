package xwall

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type marshaler struct {
	wall *Wall
}

func JSON(w *Wall) interface {
	json.Marshaler
	json.Unmarshaler
} {
	return &marshaler{wall: w}
}

func (r *marshaler) MarshalJSON() ([]byte, error) {
	r.wall.mu.RLock()
	defer r.wall.mu.RUnlock()
	b := new(bytes.Buffer)
	b.Write([]byte("["))
	fmt.Fprintf(b, `"%d"`, r.wall.lastID)
	for e := r.wall.wall.Front(); e != nil; e = e.Next() {
		fmt.Fprintf(b, `, "%s"`, e.Value.([]byte))
	}
	b.Write([]byte("]"))
	return b.Bytes(), nil
}

func (r *marshaler) UnmarshalJSON(data []byte) error {
	r.wall.mu.Lock()
	defer r.wall.mu.Unlock()
	a := []string(nil)
	err := json.Unmarshal(data, &a)
	if err != nil {
		return fmt.Errorf("wall unmarshal: %w", err)
	}
	if len(a) == 0 {
		return fmt.Errorf("wall unmarshal: unexpected empty list")
	}
	r.wall.lastID, err = strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		return fmt.Errorf("wall unmarshal: first element: %w", err)
	}
	r.wall.wall.Init()
	for _, v := range a[1:] {
		r.wall.wall.PushBack([]byte(v))
	}
	return nil
}
