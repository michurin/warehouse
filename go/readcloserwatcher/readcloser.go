package readcloserwatcher

import (
	"bytes"
	"errors"
	"io"
	"sync"
	"time"
)

type readCloserConfig struct {
	timeout time.Duration
	limit   int
}

type Result struct {
	Err    error
	Octets []byte
}

type readCloser struct {
	next   io.ReadCloser
	octets bytes.Buffer
	limit  int
	ch     chan Result
	done   chan struct{}
	lock   sync.Mutex
}

func (rc *readCloser) finalise(err error) {
	rc.lock.Lock()
	defer rc.lock.Unlock()
	if rc.ch == nil {
		return
	}
	rc.ch <- Result{
		Err:    err,
		Octets: rc.octets.Bytes(),
	}
	rc.ch = nil
	close(rc.done)

}

func (rc *readCloser) Read(p []byte) (int, error) {
	n, err := rc.next.Read(p)
	if err != nil && err != io.EOF {
		rc.finalise(err)
	}
	m := rc.limit - rc.octets.Len()
	if m <= n {
		rc.octets.Write(p[:m])
		rc.finalise(LimitError)
	} else {
		rc.octets.Write(p[:n]) // buffer can be bigger
	}
	return n, err
}

func (rc *readCloser) Close() error {
	err := rc.next.Close()
	rc.finalise(err)
	return err
}

var (
	TimeoutError = errors.New("timeout")
	LimitError   = errors.New("limit")
)

const (
	defaultTimeout = time.Minute
	defaultLimit   = 1 << 12
)

type option func(*readCloserConfig)

func Watcher(s io.ReadCloser, opts ...option) (io.ReadCloser, chan Result) {
	r := make(chan Result, 1)
	if s == nil {
		r <- Result{}
		close(r)
		return s, r
	}
	cfg := readCloserConfig{
		timeout: defaultTimeout,
		limit:   defaultLimit,
	}
	for _, o := range opts {
		o(&cfg)
	}
	done := make(chan struct{}, 1)
	timeoutChan := time.After(cfg.timeout)
	rc := &readCloser{
		limit: cfg.limit,
		next:  s,
		ch:    r,
		done:  done,
	}
	go func() {
		select {
		case <-done:
		case <-timeoutChan:
			rc.finalise(TimeoutError)
		}
	}()
	return rc, r
}

func WithTimeout(timeout time.Duration) option {
	return func(c *readCloserConfig) {
		c.timeout = timeout
	}
}

func WithLimit(limit int) option {
	return func(c *readCloserConfig) {
		c.limit = limit
	}
}
