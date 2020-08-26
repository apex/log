package journal

import (
	"errors"
	"fmt"

	"github.com/apex/log"
	"github.com/coreos/go-systemd/v22/journal"
)

// Handler implementation.
type Handler struct {
}

// New handler.
func New() *Handler {
	return &Handler{}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {

	fields := make(map[string]string, len(e.Fields))
	for k, v := range e.Fields {
		fields[k] = fmt.Sprint(v)
	}

	switch e.Level {
	case log.DebugLevel:
		return journal.Send(e.Message, journal.PriDebug, fields)
	case log.InfoLevel:
		return journal.Send(e.Message, journal.PriInfo, fields)
	case log.WarnLevel:
		return journal.Send(e.Message, journal.PriWarning, fields)
	case log.ErrorLevel:
		return journal.Send(e.Message, journal.PriErr, fields)
	case log.FatalLevel:
		return journal.Send(e.Message, journal.PriCrit, fields)
	}

	return errors.New("unknown log level given for systemd")
}

// Close shuts down the handler.
func (h *Handler) Close() error {
	return nil
}
