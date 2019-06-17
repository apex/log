// Package apexlogs implements a handler for Apex Logs.
package apexlogs

import (
	"sync"

	"github.com/apex/log"
	"github.com/apex/log/handlers/apexlogs/client"
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
	buffer []client.Event

	// client
	once sync.Once
	c    client.Client
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	// initialize client
	h.once.Do(func() {
		h.c = client.Client{
			URL: h.URL,
		}
	})

	// create event
	event := client.Event{
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

// Flush all buffered logs.
func (h *Handler) Flush() error {
	h.mu.Lock()
	events := h.buffer
	h.buffer = nil
	h.mu.Unlock()

	return h.c.AddEvents(client.AddEventsInput{
		ProjectID: h.ProjectID,
		Events:    events,
	})
}
