// Package papertrail implements a papertrail logfmt format handler.
package papertrail

import (
	"bytes"
	"fmt"
	"log/syslog"
	"net"
	"os"
	"sync"
	"time"

	"github.com/apex/log"
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

	// Conn (used for testing)
	Conn net.Conn
}

// Handler implementation.
type Handler struct {
	*Config

	mu sync.Mutex
}

// New handler.
func New(config *Config) *Handler {
	if config.Conn == nil {
		c, err := net.Dial("udp", fmt.Sprintf("%s.papertrailapp.com:%d", config.Host, config.Port))
		if err != nil {
			panic(err)
		}
		config.Conn = c
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

	h.mu.Lock()
	_, err := h.Config.Conn.Write(msg)
	h.mu.Unlock()
	return err
}

// Flush fn
func (h *Handler) Flush() {
	log.Infof("flushing!!!")
}
