package wall_test

import (
	"context"
	"encoding/json"
	"testing"

	"sse/wall"
)

func TestWall(t *testing.T) {
	// TODO split tests
	// TODO asserts
	// -- create and fill
	w := wall.New(100)
	w.Pub([]byte("message 1"))
	w.Pub([]byte("message 2"))

	// -- marshal
	b, err := json.Marshal(wall.JSON(w))
	t.Log(err)
	t.Log(string(b))

	// -- unmarshal
	w = wall.New(999)
	json.Unmarshal(b, wall.JSON(w))

	// -- check unmarshaled
	b, err = json.Marshal(wall.JSON(w))
	t.Log(err)
	t.Log(string(b))

	// -- fetch
	x, id := w.Fetch(context.TODO(), 101)
	t.Log(id)           // 102
	t.Log(len(x))       // 1
	t.Log(string(x[0])) // "message 2"
}
