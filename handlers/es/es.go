// Package es implements an Elasticsearch batch handler. Currently this implementation
// assumes the index format of "index-YY-MM-DD".
package es

import (
	"io"
	stdlog "log"
	"sync"
	"time"

	"github.com/tj/go-elastic/batch"

	"github.com/apex/log"
)

// TODO(tj): allow index configuration
// TODO(tj): allow dumping logs to stderr on timeout
// TODO(tj): allow custom format that does not include .fields etc

// index for the current time.
func index() string {
	return time.Now().Format("logs-06-01-02")
}

// Elasticsearch interface.
type Elasticsearch interface {
	Bulk(io.Reader) error
}

// Config for handler.
type Config struct {
	BufferSize int           // BufferSize is the number of logs to buffer before flush (default: 100)
	Client     Elasticsearch // Client for ES
}

// defaults applies defaults to the config.
func (c *Config) defaults() {
	if c.BufferSize == 0 {
		c.BufferSize = 100
	}
}

// Handler implementation.
type Handler struct {
	*Config

	mu    sync.Mutex
	batch *batch.Batch
}

// New handler with BufferSize
func New(config *Config) *Handler {
	config.defaults()
	return &Handler{
		Config: config,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.batch == nil {
		h.batch = &batch.Batch{
			Elastic: h.Client,
			Index:   index(),
			Type:    "log",
		}
	}

	h.batch.Add(e)

	if h.batch.Size() >= h.BufferSize {
		h.flush(h.batch)
		h.batch = nil
	}

	return nil
}

// flush the given `batch` asynchronously.
func (h *Handler) flush(batch *batch.Batch) {
	size := batch.Size()
	start := time.Now()
	stdlog.Printf("log/elastic: flushing %d logs", size)

	if err := batch.Flush(); err != nil {
		stdlog.Printf("log/elastic: failed to flush %d logs: %s", size, err)
	}

	stdlog.Printf("log/elastic: flushed %d logs in %s", size, time.Since(start))
}
