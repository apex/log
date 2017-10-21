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

	// Provide your own buffer
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
		c, err := newClient(fmt.Sprintf("%s.papertrailapp.com:%d", config.Host, config.Port))
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

func newClient(url string) (*client, error) {
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

	conn, e := net.Dial("tcp", c.url)
	if e != nil {
		return e
	}

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

func (c *client) Write(b []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		if n, err := c.conn.Write(b); err == nil {
			return n, err
		}
	}
	if err := c.connect(); err != nil {
		return 0, err
	}

	return c.conn.Write(b)
}
