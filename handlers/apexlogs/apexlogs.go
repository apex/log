// Package apexlogs implements a handler for Apex Logs.
package apexlogs

import (
	"sync"

	"github.com/apex/log"
	"github.com/apex/logs/go/logs"
)

// TODO: periodic buffering

var levelMap = map[log.Level]string{
	log.DebugLevel: "debug",
	log.InfoLevel:  "info",
	log.WarnLevel:  "warning",
	log.ErrorLevel: "error",
	log.FatalLevel: "emergency",
}

// Handler implementation.
type Handler struct {
	// URL is the endpoint for your Apex Logs deployment API.
	URL string

	// ProjectID is the id of the project that you are collecting logs for.
	ProjectID string

	// buffer
	mu     sync.Mutex
	buffer []logs.Event

	// client
	once sync.Once
	c    logs.Client
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	// initialize client
	h.once.Do(func() {
		h.c = logs.Client{
			URL: h.URL,
		}
	})

	// create event
	event := logs.Event{
		Level:     levelMap[e.Level],
		Message:   e.Message,
		Fields:    map[string]interface{}(e.Fields),
		Timestamp: e.Timestamp,
	}

	// buffer event
	h.mu.Lock()
	h.buffer = append(h.buffer, event)
	h.mu.Unlock()

	return nil
}

// Events returns the buffered events, and clears the buffer.
func (h *Handler) Events() (events []logs.Event) {
	h.mu.Lock()
	events = h.buffer
	h.buffer = nil
	h.mu.Unlock()
	return
}

// Flush all buffered logs.
func (h *Handler) Flush() error {
	return h.c.AddEvents(logs.AddEventsInput{
		ProjectID: h.ProjectID,
		Events:    h.Events(),
	})
}
