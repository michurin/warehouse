package chat

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// TODO split
	// TODO vanish sleeps

	m1 := json.RawMessage("one")
	m2 := json.RawMessage("two")
	m3 := json.RawMessage("three")
	m4 := json.RawMessage("four")

	s := newRoom(func() int64 { return 0 })
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	d, i := s.fetch(ctx, 0) // no messages, exit due to context is canceled
	assert.Nil(t, d)
	assert.Equal(t, int64(0), i)

	s.pub(m1)
	s.pub(m2)
	d, i = s.fetch(context.Background(), i)
	assert.Equal(t, []json.RawMessage{m2, m1}, d)
	assert.Equal(t, int64(2), i)

	s.pub(m3)
	d, i = s.fetch(context.Background(), i)
	assert.Equal(t, []json.RawMessage{m3}, d)
	assert.Equal(t, int64(3), i)

	go func() {
		time.Sleep(time.Microsecond)
		s.pub(m4)
	}()
	d, i = s.fetch(context.Background(), i)
	assert.Equal(t, []json.RawMessage{m4}, d)
	assert.Equal(t, int64(4), i)
}
