// Package cli implements a colored text handler suitable for command-line interfaces.
package cli

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/fatih/color"
)

// Default handler outputting to stderr.
var Default = New(os.Stderr)

// start time.
var start = time.Now()

var bold = color.New(color.Bold)

// colors.
const (
	none   = 0
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	gray   = 37
)

// Colors mapping.
var Colors = [...]*color.Color{
	log.DebugLevel: color.New(color.Attribute(gray)),
	log.InfoLevel:  color.New(color.Attribute(blue)),
	log.WarnLevel:  color.New(color.Attribute(yellow)),
	log.ErrorLevel: color.New(color.Attribute(red)),
	log.FatalLevel: color.New(color.Attribute(red)),
}

// Strings mapping.
var Strings = [...]string{
	log.DebugLevel: "•",
	log.InfoLevel:  "•",
	log.WarnLevel:  "•",
	log.ErrorLevel: "⨯",
	log.FatalLevel: "⨯",
}

// Handler implementation.
type Handler struct {
	mu      sync.Mutex
	Writer  io.Writer
	Padding int
}

// New handler.
func New(w io.Writer) *Handler {
	return &Handler{
		Writer:  w,
		Padding: 3,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	color := Colors[e.Level]
	level := Strings[e.Level]
	names := e.Fields.Names()

	h.mu.Lock()
	defer h.mu.Unlock()

	color.Fprintf(h.Writer, "%s %-25s", bold.Sprintf("%*s", h.Padding+1, level), e.Message)

	for _, name := range names {
		if name == "source" {
			continue
		}
		fmt.Fprintf(h.Writer, " %s=%s", color.Sprint(name), e.Fields.Get(name))
	}

	fmt.Fprintln(h.Writer)

	return nil
}
