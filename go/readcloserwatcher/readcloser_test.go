package readcloserwatcher_test

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"github.com/michurin/warehouse/go/readcloserwatcher"
	"github.com/stretchr/testify/assert"
)

func TestGlodenFlow(t *testing.T) {
	b := ioutil.NopCloser(bytes.NewBufferString("data"))
	w, c := readcloserwatcher.Watcher(b)
	o, err := ioutil.ReadAll(w)
	w.Close()
	r := <-c
	assert.NoError(t, err)
	assert.Equal(t, []byte("data"), o)
	assert.NoError(t, r.Err)
	assert.Equal(t, []byte("data"), r.Octets)
}

func TestNil(t *testing.T) {
	_, c := readcloserwatcher.Watcher(nil)
	r := <-c
	assert.NoError(t, r.Err)
	assert.Nil(t, r.Octets)
}

func TestLimit(t *testing.T) {
	b := ioutil.NopCloser(bytes.NewBufferString("data"))
	w, c := readcloserwatcher.Watcher(b, readcloserwatcher.WithLimit(2))
	o, err := ioutil.ReadAll(w)
	w.Close()
	r := <-c
	assert.NoError(t, err)
	assert.Equal(t, []byte("data"), o)
	assert.EqualError(t, r.Err, "limit")
	assert.Equal(t, []byte("da"), r.Octets)
}

func TestTimeout_naive(t *testing.T) {
	b := ioutil.NopCloser(bytes.NewBufferString("data"))
	w, c := readcloserwatcher.Watcher(b, readcloserwatcher.WithTimeout(50*time.Millisecond))
	time.Sleep(100 * time.Millisecond)
	o, err := ioutil.ReadAll(w)
	w.Close()
	r := <-c
	assert.NoError(t, err)
	assert.Equal(t, []byte("data"), o)
	assert.EqualError(t, r.Err, "timeout")
	assert.Equal(t, []byte(nil), r.Octets)
}

func TestTimeout2_naive(t *testing.T) {
	b := ioutil.NopCloser(bytes.NewBufferString("data"))
	w, c := readcloserwatcher.Watcher(b, readcloserwatcher.WithTimeout(50*time.Millisecond))
	o, err := ioutil.ReadAll(w)
	time.Sleep(100 * time.Millisecond)
	w.Close()
	r := <-c
	assert.NoError(t, err)
	assert.Equal(t, []byte("data"), o)
	assert.EqualError(t, r.Err, "timeout")
	assert.Equal(t, []byte("data"), r.Octets)
}
