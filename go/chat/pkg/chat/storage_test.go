package chat

import (
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	s := New()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	d, i := s.Get(ctx, 0)
	t.Log(d)
	t.Log(i)
	s.Add(Message{Text: "one"})
	s.Add(Message{Text: "two"})
	d, i = s.Get(context.Background(), i)
	t.Log(d)
	t.Log(i)
	s.Add(Message{Text: "three"})
	d, i = s.Get(context.Background(), i)
	t.Log(d)
	t.Log(i)
	go func() {
		time.Sleep(time.Microsecond)
		s.Add(Message{Text: "four"})
	}()
	d, i = s.Get(context.Background(), i)
	t.Log(d)
	t.Log(i)
}
