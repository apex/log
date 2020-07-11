// Package apexlogs implements a handler for Apex Logs https://apex.sh/logs/.
package apexlogs

import (
	"context"
	stdlog "log"
	"net/http"
	"os"
	"sync"

	"github.com/tj/go-buffer"

	"github.com/apex/log"
	"github.com/apex/logs/go/logs"
)

// logger instance.
var logger = stdlog.New(os.Stderr, "buffer ", stdlog.LstdFlags)

// levelMap is a mapping of severity levels.
var levelMap = map[log.Level]string{
	log.DebugLevel: "debug",
	log.InfoLevel:  "info",
	log.WarnLevel:  "warning",
	log.ErrorLevel: "error",
	log.FatalLevel: "emergency",
}

// Handler implementation.
type Handler struct {
	url           string
	projectID     string
	authToken     string
	httpClient    *http.Client
	bufferOptions []buffer.Option

	once sync.Once
	b    *buffer.Buffer
	c    logs.Client
}

// Option function.
type Option func(*Handler)

// New Apex Logs handler with the url, projectID and options.
func New(url, projectID string, options ...Option) *Handler {
	var v Handler
	v.url = url
	v.projectID = projectID
	for _, o := range options {
		o(&v)
	}
	return &v
}

// WithAuthToken sets the authentication token used for requests.
func WithAuthToken(token string) Option {
	return func(v *Handler) {
		v.authToken = token
	}
}

// WithHTTPClient sets the HTTP client used for requests.
func WithHTTPClient(client *http.Client) Option {
	return func(v *Handler) {
		v.httpClient = client
	}
}

// WithBufferOptions sets options for the underlying buffer used to batch logs.
func WithBufferOptions(options ...buffer.Option) Option {
	return func(v *Handler) {
		v.bufferOptions = options
	}
}

// init the client and buffer.
func (h *Handler) init() {
	h.c = logs.Client{
		URL:        h.url,
		HTTPClient: h.httpClient,
		AuthToken:  h.authToken,
	}

	var options []buffer.Option
	options = append(options, buffer.WithFlushHandler(h.handleFlush))
	options = append(options, buffer.WithErrorHandler(h.handleError))
	options = append(options, h.bufferOptions...)
	h.b = buffer.New(options...)
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	h.once.Do(h.init)

	h.b.Push(logs.Event{
		Level:     levelMap[e.Level],
		Message:   e.Message,
		Fields:    map[string]interface{}(e.Fields),
		Timestamp: e.Timestamp,
	})

	return nil
}

// Flush any pending logs.
func (h *Handler) Flush() {
	h.b.Flush()
}

// Close flushes any pending logs, and waits for flushing to complete. This
// method should be called before exiting your program to ensure entries have
// flushed properly.
func (h *Handler) Close() {
	h.b.Close()
}

// handleFlush implementation.
func (h *Handler) handleFlush(ctx context.Context, values []interface{}) error {
	var events []logs.Event

	for _, v := range values {
		events = append(events, v.(logs.Event))
	}

	return h.c.AddEvents(logs.AddEventsInput{
		ProjectID: h.projectID,
		Events:    events,
	})
}

// handleError implementation.
func (h *Handler) handleError(err error) {
	logger.Printf("error flushing logs: %v", err)
}
