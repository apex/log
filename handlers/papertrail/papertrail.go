// Package papertrail implements a papertrail logfmt format handler.
package papertrail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	stdlog "log"
	"log/syslog"
	"net"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/internal/queue"
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

	// Advanced: tweak how we send logs to papertrail
	// Capacity is the amount of log to queue before
	// we start dropping, write timeout dictates how
	// long we should wait before giving up on a write
	Capacity       int
	ConnectTimeout time.Duration
	WriteTimeout   time.Duration
}

// Handler implementation.
type Handler struct {
	url  string
	c    *Config
	conn net.Conn
	q    *queue.Queue
}

// New handler.
func New(config *Config) *Handler {

	// defaults
	// TODO: what should these be?
	if config.Capacity == 0 {
		config.Capacity = 20
	}
	if config.ConnectTimeout == 0 {
		config.ConnectTimeout = 30 * time.Second
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 5 * time.Second
	}

	// connect to papertrail
	url := fmt.Sprintf("%s.papertrailapp.com:%d", config.Host, config.Port)
	conn, err := connect(url, config.ConnectTimeout)
	if err != nil {
		// TODO: should we kill here?
		stdlog.Printf("log/papertrail: couldn't connect to papertrail")
	}

	// TODO: see if papertrail can handle out of order
	// logs. if so, we can make our write function
	// thread-safe & adjust the concurrency here
	q := queue.New(config.Capacity, 1)

	return &Handler{
		c:    config,
		conn: conn,
		url:  url,
		q:    q,
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

	msg := []byte(fmt.Sprintf("<%d>%s %s %s[%d]: %s\n",
		syslog.LOG_KERN,
		ts,
		h.c.Hostname,
		h.c.Tag,
		os.Getpid(),
		buf.String(),
	))

	return h.q.Push(func() {
		if e := h.write(msg); e != nil {
			stdlog.Printf("log/papertrail: %s", e)
		}
	})
}

// connect to papertrail
// based on: https://github.com/papertrail/remote_syslog2/blob/master/syslog/syslog.go
func connect(url string, connectTimeout time.Duration) (net.Conn, error) {
	config := &tls.Config{
		RootCAs: certpool(),
	}
	dialer := &net.Dialer{
		Timeout:   connectTimeout,
		KeepAlive: time.Second * 60 * 3, // 3 minutes
	}
	return tls.DialWithDialer(dialer, "tcp", url, config)
}

// write the data to the connection
// this is not thread-safe
func (h *Handler) write(b []byte) error {
	deadline := time.Now().Add(h.c.WriteTimeout)

	// set a connection deadline, so writes don't
	// block forever
	if e := h.conn.SetWriteDeadline(deadline); e != nil {
		return e
	}

	// try writing on the existing connection
	n, e := h.conn.Write(b)
	if e == nil {
		return nil
	}

	// try reconnecting in the event of an error
	c, e := connect(h.url, h.c.ConnectTimeout)
	if e != nil {
		return e
	}
	h.conn = c

	// try writing again if we made a successful connection
	_, e = h.conn.Write(b[n:])
	if e != nil {
		return e
	}

	return nil
}

// Flush fn
func (h *Handler) Flush() {
	h.q.Wait()
}
