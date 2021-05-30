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
	s := New(0)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	d, i := s.Get(ctx, 0) // no messages, exit due to context is canceled
	assert.Nil(t, d)
	assert.Equal(t, 0, i)
	s.Put(json.RawMessage("one"))
	s.Put(json.RawMessage("two"))
	d, i = s.Get(context.Background(), i)
	assert.Equal(t, []Message{
		{Message: json.RawMessage("two")},
		{Message: json.RawMessage("one")},
	}, d)
	assert.Equal(t, 2, i)
	s.Put(json.RawMessage("three"))
	d, i = s.Get(context.Background(), i)
	assert.Equal(t, []Message{
		{Message: json.RawMessage("three")},
	}, d)
	assert.Equal(t, 3, i)
	go func() {
		time.Sleep(time.Microsecond)
		s.Put(json.RawMessage("four"))
	}()
	d, i = s.Get(context.Background(), i)
	assert.Equal(t, []Message{
		{Message: json.RawMessage("four")},
	}, d)
	assert.Equal(t, 4, i)
}
