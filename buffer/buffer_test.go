package buffer_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/apex/log/buffer"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	var w bytes.Buffer
	buf := buffer.NewWithConfig(&w, &buffer.Config{
		BufferSize: 1,
	})
	buf.Append([]byte("a"))
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, "a", w.String())
}

func TestOverCap(t *testing.T) {
	var w bytes.Buffer
	buf := buffer.NewWithConfig(&w, &buffer.Config{
		BufferSize: 1,
	})
	buf.Append([]byte("aa"))
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, "aa", w.String())
}

func TestMultiple(t *testing.T) {
	var w bytes.Buffer
	buf := buffer.NewWithConfig(&w, &buffer.Config{
		BufferSize: 2,
	})
	buf.Append([]byte("a"))
	buf.Append([]byte("aa"))
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, "aaa", w.String())
}
func TestMultiple2(t *testing.T) {
	var w bytes.Buffer
	buf := buffer.NewWithConfig(&w, &buffer.Config{
		BufferSize: 2,
	})
	buf.Append([]byte("aa"))
	buf.Append([]byte("a"))
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, "aa", w.String())
}

func TestHalfFull(t *testing.T) {
	var w bytes.Buffer
	buf := buffer.NewWithConfig(&w, &buffer.Config{
		BufferSize: 2,
	})
	buf.Append([]byte("a"))
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, "", w.String())
}

func TestHalfFullFlush(t *testing.T) {
	var w bytes.Buffer
	buf := buffer.NewWithConfig(&w, &buffer.Config{
		BufferSize: 2,
	})
	buf.Append([]byte("a"))
	buf.Flush()
	assert.Equal(t, "a", w.String())
}

func TestIdleFlush(t *testing.T) {
	var w bytes.Buffer
	buf := buffer.NewWithConfig(&w, &buffer.Config{
		BufferSize:  2,
		IdleTimeout: 50 * time.Millisecond,
	})
	buf.Append([]byte("a"))
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, "a", w.String())
}

type SlowWriter struct {
	buf bytes.Buffer
}

func (s *SlowWriter) Write(b []byte) (int, error) {
	time.Sleep(500 * time.Millisecond)
	return s.buf.Write(b)
}

func (s *SlowWriter) String() string {
	return s.buf.String()
}

func TestDroppedWrites(t *testing.T) {
	w := SlowWriter{}
	buf := buffer.NewWithConfig(&w, &buffer.Config{
		BufferSize:     1,
		PendingFlushes: 1,
	})
	buf.Append([]byte("a"))
	buf.Append([]byte("b"))
	buf.Append([]byte("c"))
	buf.Flush()
	assert.Equal(t, "ab", w.String())
}

func TestDroppedWrites2(t *testing.T) {
	w := SlowWriter{}
	buf := buffer.NewWithConfig(&w, &buffer.Config{
		BufferSize:     1,
		PendingFlushes: 1,
	})
	buf.Append([]byte("abc"))
	buf.Append([]byte("def"))
	buf.Append([]byte("ghi"))
	buf.Flush()
	assert.Equal(t, "abcdef", w.String())
}

type UnavailableWriter struct {
	buf bytes.Buffer
}

func (s *UnavailableWriter) Write(b []byte) (int, error) {
	return 0, errors.New("unavailable")
}

func (s *UnavailableWriter) String() string {
	return s.buf.String()
}

func TestUnavailable(t *testing.T) {
	w := UnavailableWriter{}
	buf := buffer.NewWithConfig(&w, &buffer.Config{
		BufferSize:     1,
		PendingFlushes: 1,
	})
	buf.Append([]byte("abc"))
	buf.Append([]byte("def"))
	buf.Append([]byte("ghi"))
	buf.Append([]byte("ghi"))
	buf.Append([]byte("ghi"))
	buf.Append([]byte("ghi"))
	buf.Append([]byte("ghi"))
	buf.Append([]byte("ghi"))
	buf.Append([]byte("ghi"))
	buf.Append([]byte("ghi"))
	buf.Append([]byte("ghi"))
	buf.Append([]byte("ghi"))
	buf.Append([]byte("ghi"))
	buf.Flush()
	assert.Equal(t, "", w.String())
}

type FlakyWriter struct {
	n   int
	buf bytes.Buffer
}

func (s *FlakyWriter) Write(b []byte) (int, error) {
	s.n++
	if s.n%2 == 0 {
		return 0, errors.New("unavailable")
	} else {
		return s.buf.Write(b)
	}
}

func (s *FlakyWriter) String() string {
	return s.buf.String()
}
func TestFlaky(t *testing.T) {
	w := FlakyWriter{}
	buf := buffer.NewWithConfig(&w, &buffer.Config{
		BufferSize: 1,
	})
	buf.Append([]byte("abc"))
	buf.Append([]byte("def"))
	buf.Append([]byte("ghi"))
	buf.Flush()
	assert.Equal(t, "abcdefghi", w.String())
}

type LeakyWriter struct {
	n   int
	buf bytes.Buffer
}

func (s *LeakyWriter) Write(b []byte) (int, error) {
	s.n++

	n := len(b)
	if s.n%2 == 1 {
		n = len(b) / 2
	}

	s.buf.Write(b[:n])
	return len(b[:n]), nil
}

func (s *LeakyWriter) String() string {
	return s.buf.String()
}
func TestLeaky(t *testing.T) {
	w := LeakyWriter{}
	buf := buffer.New(&w)
	buf.Append([]byte("ab"))
	buf.Append([]byte("cd"))
	buf.Append([]byte("ef"))
	buf.Flush()
	assert.Equal(t, "abcdef", w.String())
}
