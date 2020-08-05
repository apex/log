// Package apexlogs implements a handler for Apex Logs https://apex.sh/logs/.
package apexlogs

import (
	"context"
	stdlog "log"
	"net/http"
	"os"

	"github.com/tj/go-buffer"

	"github.com/apex/log"
	"github.com/apex/logs"
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
	projectID     string
	httpClient    *http.Client
	bufferOptions []buffer.Option

	b *buffer.Buffer
	c logs.Client
}

// Option function.
type Option func(*Handler)

// New Apex Logs handler with the url, projectID, authToken and options.
func New(url, projectID, authToken string, options ...Option) *Handler {
	var v Handler
	v.projectID = projectID

	// options
	for _, o := range options {
		o(&v)
	}

	// logs client
	v.c = logs.Client{
		URL:        url,
		AuthToken:  authToken,
		HTTPClient: v.httpClient,
	}

	// event buffer
	var o []buffer.Option
	o = append(o, buffer.WithFlushHandler(v.handleFlush))
	o = append(o, buffer.WithErrorHandler(v.handleError))
	o = append(o, v.bufferOptions...)
	v.b = buffer.New(o...)

	return &v
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

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	h.b.Push(logs.Event{
		Level:     levelMap[e.Level],
		Message:   e.Message,
		Fields:    map[string]interface{}(e.Fields),
		Timestamp: e.Timestamp,
	})

	return nil
}

// Flush any pending logs. This method is non-blocking.
func (h *Handler) Flush() {
	h.b.Flush()
}

// FlushSync any pending logs. This method is blocking.
func (h *Handler) FlushSync() {
	h.b.FlushSync()
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

	if len(events) == 0 {
		return nil
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
