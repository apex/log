package buffer

import (
	"io"
	stdlog "log"
	"time"
)

// Buffer struct
type Buffer struct {
	w io.Writer
	c *Config

	cmdc   chan func()
	writec chan *request

	buf []byte
}

// Config struct
type Config struct {
	// Buffer size in bytes before flushing
	BufferSize int
	// Time to wait without an buf.Append() before flushing
	IdleTimeout time.Duration
	// Number of flushes to queue up before we start dropping
	PendingFlushes int
	// Number of time to retry a write before giving up
	RetryWrites int
}

// write request with an optional drain parameter
type request struct {
	buf    []byte
	drainc chan bool
}

// New buffer. This buffer never blocks on writes but may drop
// logs if the receiver is either down or too slow. It will also
// flush during idle periods
func New(writer io.Writer) *Buffer {
	return NewWithConfig(writer, &Config{
		BufferSize:     1028,
		IdleTimeout:    5 * time.Second,
		PendingFlushes: 100,
		RetryWrites:    3,
	})
}

// NewWithConfig gives us more advanced options to configure
func NewWithConfig(writer io.Writer, config *Config) *Buffer {
	if writer == nil {
		panic("buffer needs a writer")
	}
	if config.BufferSize == 0 {
		config.BufferSize = 1028
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = 5 * time.Second
	}
	if config.PendingFlushes == 0 {
		config.PendingFlushes = 100
	}
	if config.RetryWrites == 0 {
		config.RetryWrites = 3
	}

	cmdc := make(chan func(), 1)
	writec := make(chan *request, config.PendingFlushes)

	b := &Buffer{
		w: writer,
		c: config,

		cmdc:   cmdc,
		writec: writec,
	}

	go b.loop()
	go b.write()

	return b
}

// event loop for commands that flushes after being
// idle for the specified duration
func (b *Buffer) loop() {
	for {
		select {
		case fn := <-b.cmdc:
			fn()
		case <-time.After(b.c.IdleTimeout):
			b.idleFlush()
		}
	}
}

// write with retries
func (b *Buffer) write() {
	for {
		req := <-b.writec
		buf := req.buf

		// retry the write a few times
		// before giving up & moving on
		retries := b.c.RetryWrites
		for len(buf) > 0 {
			n, e := b.w.Write(buf)

			// done writing
			if e == nil && len(buf) == n {
				break
			}

			buf = buf[n:]
			retries--

			if retries <= 0 {
				stdlog.Printf("log/buffer: error writing buffer '%s'", buf)
				break
			}
		}

		// ack that we've drained
		// if the requester include
		// this channel
		if req.drainc != nil {
			close(req.drainc)
		}
	}
}

// Append to the buffer. This will periodically flush
// depending on the buffer size. This will never block.
func (b *Buffer) Append(msg []byte) {
	b.cmdc <- func() {
		b.buf = append(b.buf, msg...)
		if len(b.buf) >= b.c.BufferSize {
			d := make([]byte, len(b.buf))
			copy(d, b.buf)
			b.buf = b.buf[:0]
			b.flush(d)
		}
	}
}

// Flush the buffer manually. This will block until
// the writes have finished. This will also block
// writes following this call so you'll probably
// want to call this at the end
func (b *Buffer) Flush() {
	drainc := make(chan bool)
	b.cmdc <- func() {
		d := make([]byte, len(b.buf))
		copy(d, b.buf)
		b.buf = b.buf[:0]
		b.writec <- &request{d, drainc}
	}
	<-drainc
}

// non-blocking manual flush
func (b *Buffer) idleFlush() {
	b.cmdc <- func() {
		if len(b.buf) == 0 {
			return
		}
		d := make([]byte, len(b.buf))
		copy(d, b.buf)
		b.flush(d)
	}
}

// write the buffer, dropping new requests
// if there is too much backpressure
func (b *Buffer) flush(buf []byte) {
	select {
	case b.writec <- &request{buf: buf}:
	default:
		// drop the request
	}
}
