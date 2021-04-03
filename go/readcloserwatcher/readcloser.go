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
}

type result struct {
	Err    error
	Octets []byte
}

type readCloser struct {
	next   io.ReadCloser
	octets bytes.Buffer
	ch     chan result
	done   chan struct{}
	lock   sync.Mutex
}

func (rc *readCloser) finalise(err error) {
	rc.lock.Lock()
	defer rc.lock.Unlock()
	if rc.ch == nil {
		return
	}
	rc.ch <- result{
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
	rc.octets.Write(p[:n]) // buffer can be bigger
	return n, err
}

func (rc *readCloser) Close() error {
	err := rc.next.Close()
	rc.finalise(err)
	return err
}

var TimeoutError = errors.New("timeout")

const defaultTimeout = time.Minute

type option func(*readCloserConfig)

func Watcher(s io.ReadCloser, opts ...option) (io.ReadCloser, chan result) {
	r := make(chan result, 1)
	if s == nil {
		r <- result{}
		close(r)
		return s, r
	}
	cfg := readCloserConfig{
		timeout: defaultTimeout,
	}
	for _, o := range opts {
		o(&cfg)
	}
	done := make(chan struct{}, 1)
	timeoutChan := time.After(cfg.timeout)
	rc := &readCloser{
		next: s,
		ch:   r,
		done: done,
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
