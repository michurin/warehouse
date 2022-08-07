package stream_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/michurin/warehouse/go/chat2/stream"
)

func assertGet(t *testing.T, expBytes [][]byte, expBound uint64, actBytes [][]byte, actBound uint64) {
	t.Helper()
	if expBound != actBound {
		t.Errorf("Bound: exp: %d, but got: %d", expBound, actBound)
	}
	if (expBytes == nil && actBytes != nil) || (expBytes != nil && actBytes == nil) {
		t.Errorf("Bytes: exp: %v, but got: %v", expBytes, actBytes)
	}
	if len(expBytes) != len(actBytes) {
		t.Errorf("Bytes: len: exp: %v, but got: %v", expBytes, actBytes)
	}
	for i, v := range expBytes {
		if !bytes.Equal(actBytes[i], v) {
			t.Errorf("Exp: %v, but got: %v (i=%d)", expBytes, actBytes, i)
		}
	}
}

func bt(s string) []byte {
	return []byte(s)
}

func st(s ...string) [][]byte {
	r := make([][]byte, len(s))
	for i, v := range s {
		r[i] = []byte(v)
	}
	return r
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

func TestWithouWaiting(t *testing.T) {
	for _, tt := range []struct {
		name          string
		init          [][]byte
		bound         uint64
		expMsg        [][]byte
		expContinuity bool
	}{
		{
			name:          "newReader",
			init:          st("one", "two"),
			bound:         0,
			expMsg:        st("one", "two"),
			expContinuity: true,
		}, {
			name:          "readerFromFuture", // after service restart
			init:          st("one", "two"),
			bound:         9,
			expMsg:        st("one", "two"),
			expContinuity: false,
		}, {
			name:          "readTail", // #0 has been already read
			init:          st("one", "two"),
			bound:         1,
			expMsg:        st("two"),
			expContinuity: true,
		}, {
			name:          "newReaderAll",
			init:          st("one", "two", "three", "four"),
			bound:         0,
			expMsg:        st("two", "three", "four"),
			expContinuity: true,
		}, {
			name:          "readTailAll",
			init:          st("one", "two", "three", "four"),
			bound:         1,
			expMsg:        st("two", "three", "four"),
			expContinuity: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := stream.New(3)
			for _, m := range tt.init {
				s.Put(m)
			}
			ctx := context.Background()
			a, b := s.Get(ctx, tt.bound)
			assertGet(t, tt.expMsg, uint64(len(tt.init)), a, b)
		})
	}
}

func TestWithTimeout(t *testing.T) {
	for _, tt := range []struct {
		name string
		init [][]byte
	}{{
		name: "newServeNewClient",
		init: nil,
	}, {
		name: "justWait",
		init: st("one"),
	}} {
		t.Run(tt.name, func(t *testing.T) {
			s := stream.New(3)
			for _, m := range tt.init {
				s.Put(m)
			}
			fin := make(chan struct{})
			done := make(chan struct{})
			doneRes := make(chan struct{})
			close(doneRes)
			ctx := &fakeContext{t: t, done: done, doneRes: doneRes}
			go func() {
				a, b := s.Get(ctx, uint64(len(tt.init)))
				assertGet(t, nil, uint64(len(tt.init)), a, b)
				close(fin)
			}()
			<-done // just to be sure, we are reach ctx.Done() call
			<-fin  // to be sure all asserts are done
		})
	}
}

func TestWithWaiting(t *testing.T) {
	s := stream.New(3)
	s.Put(bt("one"))
	a, b := s.Get(context.Background(), 0)
	assertGet(t, st("one"), 1, a, b)
	fin := make(chan struct{})
	done := make(chan struct{})
	ctx := &fakeContext{t: t, done: done}
	go func() {
		a, b := s.Get(ctx, b)
		assertGet(t, st("two"), 2, a, b)
		close(fin)
	}()
	<-done // make Put after Get stars waiting
	s.Put(bt("two"))
	<-fin // waiting for all assertions will be complete
}

func TestDataCurruption(t *testing.T) {
	ctx := context.Background()
	s := stream.New(3)
	s.Put(bt("one"))
	s.Put(bt("two"))
	a, b := s.Get(ctx, 0)
	assertGet(t, st("one", "two"), 2, a, b)
	s.Put(bt("3"))                          // running out ring buffer capacity
	s.Put(bt("1"))                          // overwrite "one"
	assertGet(t, st("one", "two"), 2, a, b) // a is not corrupted
	a, b = s.Get(ctx, 0)
	assertGet(t, st("two", "3", "1"), 4, a, b)
}
