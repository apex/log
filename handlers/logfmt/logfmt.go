// Package logfmt implements a "logfmt" format handler.
package logfmt

import (
	"io"
	"os"
	"sort"
	"sync"

	"github.com/apex/log"
	"github.com/go-logfmt/logfmt"
)

// Default handler outputting to stderr.
var Default = New(os.Stderr)

// field used for sorting.
type field struct {
	Name  string
	Value interface{}
}

// by sorts fields by name.
type byName []field

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// Handler implementation.
type Handler struct {
	mu  sync.Mutex
	enc *logfmt.Encoder
}

// New handler.
func New(w io.Writer) *Handler {
	return &Handler{
		enc: logfmt.NewEncoder(w),
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	var fields []field

	for k, v := range e.Fields {
		fields = append(fields, field{k, v})
	}

	sort.Sort(byName(fields))

	h.mu.Lock()
	defer h.mu.Unlock()

	h.enc.EncodeKeyval("timestamp", e.Timestamp)
	h.enc.EncodeKeyval("level", e.Level.String())
	h.enc.EncodeKeyval("message", e.Message)

	for _, f := range fields {
		h.enc.EncodeKeyval(f.Name, f.Value)
	}

	h.enc.EndRecord()

	return nil
}
