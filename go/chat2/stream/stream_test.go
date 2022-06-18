package stream_test

import (
	"context"
	"testing"
	"time"

	"github.com/michurin/warehouse/go/chat2/stream"
)

func assertInt(t *testing.T, exp, act int) {
	t.Helper()
	if exp != act {
		t.Errorf("Exp: %d, but got: %d", exp, act)
	}
}

func assertStrSlice(t *testing.T, exp, act []string) {
	t.Helper()
	if len(exp) != len(act) {
		t.Errorf("Exp: %v, but got: %v", exp, act)
	}
	for i, v := range exp {
		if act[i] != v {
			t.Errorf("Exp: %v, but got: %v (i=%d)", exp, act, i)
		}
	}
}

func assertTrue(t *testing.T, b bool) {
	t.Helper()
	if !b {
		t.Error("True expected")
	}
}

type fakeContext struct {
	t       *testing.T
	done    chan struct{}
	doneRes chan struct{} // nil will wait forever
}

func (fc *fakeContext) Deadline() (deadline time.Time, ok bool) {
	fc.t.Error("Deadline called")
	return time.Time{}, false
}

func (fc *fakeContext) Done() <-chan struct{} {
	close(fc.done)
	return fc.doneRes
}

func (fc *fakeContext) Err() error {
	fc.t.Error("Err called")
	return nil
}

func (fc *fakeContext) Value(key interface{}) interface{} {
	fc.t.Error("Value called")
	return nil
}

func byteToString(a [][]byte) []string {
	b := make([]string, len(a))
	for i, v := range a {
		b[i] = string(v)
	}
	return b
}

func TestStreamWithouWaiting(t *testing.T) {
	for _, tt := range []struct {
		name   string
		init   []string
		bound  int
		expMsg []string
	}{{
		name:   "newReader",
		init:   []string{"one", "two"},
		bound:  -1,
		expMsg: []string{"two", "one"},
	}, {
		name:   "readerFromFuture", // after service restart
		init:   []string{"one", "two"},
		bound:  2,
		expMsg: []string{"two", "one"},
	}, {
		name:   "readTail", // #0 has been already read
		init:   []string{"one", "two"},
		bound:  0,
		expMsg: []string{"two"},
	}, {
		name:   "newReaderAll",
		init:   []string{"one", "two", "three", "four"},
		bound:  -1,
		expMsg: []string{"four", "three", "two"},
	}, {
		name:   "readerFromFutureAll",
		init:   []string{"one", "two", "three", "four"},
		bound:  1,
		expMsg: []string{"four", "three"},
	}, {
		name:   "readTailAll",
		init:   []string{"one", "two", "three", "four"},
		bound:  0,
		expMsg: []string{"four", "three", "two"},
	}} {
		t.Run(tt.name, func(t *testing.T) {
			s := stream.New(3)
			for _, m := range tt.init {
				s.Put([]byte(m))
			}
			ctx := context.Background()
			a, b := s.Get(ctx, tt.bound)
			assertStrSlice(t, tt.expMsg, byteToString(a))
			assertInt(t, len(tt.init)-1, b)
		})
	}
}

func TestWithTimeout(t *testing.T) {
	for _, tt := range []struct {
		name string
		init []string
	}{{
		name: "newServeNewClient",
		init: nil,
	}, {
		name: "justWait",
		init: []string{"one"},
	}} {
		t.Run(tt.name, func(t *testing.T) {
			s := stream.New(3)
			for _, m := range tt.init {
				s.Put([]byte(m))
			}
			fin := make(chan struct{})
			done := make(chan struct{})
			doneRes := make(chan struct{})
			close(doneRes)
			ctx := &fakeContext{t: t, done: done, doneRes: doneRes}
			go func() {
				a, b := s.Get(ctx, len(tt.init)-1)
				assertTrue(t, a == nil)
				assertInt(t, len(tt.init)-1, b)
				close(fin)
			}()
			<-done // just to be sure, we are reach ctx.Done() call
			<-fin  // to be sure all asserts are done
		})
	}
}

func TestWithWaiting(t *testing.T) {
	s := stream.New(3)
	s.Put([]byte("one"))
	a, b := s.Get(context.Background(), -1)
	assertStrSlice(t, []string{"one"}, byteToString(a))
	assertInt(t, 0, b)
	fin := make(chan struct{})
	done := make(chan struct{})
	ctx := &fakeContext{t: t, done: done}
	go func() {
		a, b := s.Get(ctx, b)
		assertStrSlice(t, []string{"two"}, byteToString(a))
		assertInt(t, 1, b)
		close(fin)
	}()
	<-done // make Put after Get stars waiting
	s.Put([]byte("two"))
	<-fin // waiting for all assertions will be complete
}
