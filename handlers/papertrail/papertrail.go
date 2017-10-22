// Package papertrail implements a papertrail logfmt format handler.
package papertrail

import (
	"bytes"
	"fmt"
	"io"
	"log/syslog"
	"net"
	"os"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/buffer"
	"github.com/go-logfmt/logfmt"
)

// TODO: syslog portion is ad-hoc for my serverless use-case,
// I don't really need hostnames etc, but this should be improved

// Config for Papertrail.
type Config struct {
	// Papertrail settings.
	Host string // Host subdomain such as "logs4"
	Port int    // Port number

	// Application settings
	Hostname string // Hostname value
	Tag      string // Tag value

	// Optionally provide your own buffer
	Buffer *buffer.Buffer

	// Provide your own writer (useful for testing)
	Writer io.Writer
}

// Handler implementation.
type Handler struct {
	*Config
}

// New handler.
func New(config *Config) *Handler {
	if config.Writer == nil {
		c, err := Client(fmt.Sprintf("%s.papertrailapp.com:%d", config.Host, config.Port))
		if err != nil {
			panic(err)
		}
		config.Writer = c
	}

	if config.Buffer == nil {
		config.Buffer = buffer.New(config.Writer)
	}

	return &Handler{
		Config: config,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	ts := log.Now().Format(time.Stamp)
	var buf bytes.Buffer

	enc := logfmt.NewEncoder(&buf)
	enc.EncodeKeyval("level", e.Level.String())
	enc.EncodeKeyval("message", e.Message)

	for k, v := range e.Fields {
		enc.EncodeKeyval(k, v)
	}

	enc.EndRecord()

	msg := []byte(fmt.Sprintf("<%d>%s %s %s[%d]: %s\n", syslog.LOG_KERN, ts, h.Hostname, h.Tag, os.Getpid(), buf.String()))
	h.Buffer.Append(msg)
	return nil
}

// Flush fn
func (h *Handler) Flush() {
	h.Buffer.Flush()
}

type client struct {
	url  string
	conn net.Conn
	mu   sync.Mutex
}

// This client is based on the syslog client here:
// https://golang.org/src/log/syslog/syslog.go
//
// TODO: should this be a long-lived TCP connection?
func Client(url string) (*client, error) {
	c := &client{url: url}

	// make an initial connection (this will be long-lived)
	if e := c.connect(); e != nil {
		return nil, e
	}

	return c, nil
}

func (c *client) connect() error {
	if c.conn != nil {
		// ignore err from close, it makes sense to continue anyway
		c.conn.Close()
		c.conn = nil
	}

	addr, e := net.ResolveTCPAddr("tcp", c.url)
	if e != nil {
		return e
	}

	conn, e := net.DialTCP("tcp", nil, addr)
	if e != nil {
		return e
	}
	if e := conn.SetWriteBuffer(1); e != nil {
		return e
	}
	if e := conn.SetReadBuffer(1); e != nil {
		return e
	}
	conn.SetNoDelay(true)
	c.conn = conn
	return nil
}

func (c *client) close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}

// Writer that will first try writing, then
// if that fails, try reconnecting the TCP
// connection and writing again.

// Failures will be retried by the buffer,
// so we need to keep track of what's already
// been written.
func (c *client) Write(b []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var written int
	if c.conn != nil {
		if n, err := c.conn.Write(b); err == nil {
			fmt.Printf("n written %d\n", n)
			return n, nil
		} else if n > 0 {
			written = n
		}
	}
	if err := c.connect(); err != nil {
		return written, err
	}

	return c.conn.Write(b[written:])
}
