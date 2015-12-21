package log_test

import (
	"fmt"
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/apex/log/handlers/memory"
	"github.com/stretchr/testify/assert"
)

func TestLogger_printf(t *testing.T) {
	h := memory.New()

	l := &log.Logger{
		Handler: h,
		Level:   log.InfoLevel,
	}

	l.Infof("logged in %s", "Tobi")

	e := h.Entries[0]
	assert.Equal(t, e.Message, "logged in Tobi")
	assert.Equal(t, e.Level, log.InfoLevel)
}

func TestLogger_levels(t *testing.T) {
	h := memory.New()

	l := &log.Logger{
		Handler: h,
		Level:   log.InfoLevel,
	}

	l.Debug("uploading")
	l.Info("upload complete")

	assert.Equal(t, 1, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, e.Message, "upload complete")
	assert.Equal(t, e.Level, log.InfoLevel)
}

func TestLogger_WithFields(t *testing.T) {
	h := memory.New()

	l := &log.Logger{
		Handler: h,
		Level:   log.InfoLevel,
	}

	ctx := l.WithFields(log.Fields{"file": "sloth.png"})
	ctx.Debug("uploading")
	ctx.Info("upload complete")

	assert.Equal(t, 1, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, e.Message, "upload complete")
	assert.Equal(t, e.Level, log.InfoLevel)
	assert.Equal(t, log.Fields{"file": "sloth.png"}, e.Fields)
}

func TestLogger_WithField(t *testing.T) {
	h := memory.New()

	l := &log.Logger{
		Handler: h,
		Level:   log.InfoLevel,
	}

	ctx := l.WithField("file", "sloth.png").WithField("user", "Tobi")
	ctx.Debug("uploading")
	ctx.Info("upload complete")

	assert.Equal(t, 1, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, e.Message, "upload complete")
	assert.Equal(t, e.Level, log.InfoLevel)
	assert.Equal(t, log.Fields{"file": "sloth.png", "user": "Tobi"}, e.Fields)
}

func TestLogger_Trace_info(t *testing.T) {
	h := memory.New()

	l := &log.Logger{
		Handler: h,
		Level:   log.InfoLevel,
	}

	trace := l.WithField("file", "sloth.png").Trace("upload")
	trace.Stop(nil)

	assert.Equal(t, 2, len(h.Entries))

	{
		e := h.Entries[0]
		assert.Equal(t, e.Message, "upload")
		assert.Equal(t, e.Level, log.InfoLevel)
		assert.Equal(t, log.Fields{"file": "sloth.png", "complete": false}, e.Fields)
	}

	{
		e := h.Entries[1]
		assert.Equal(t, e.Message, "upload")
		assert.Equal(t, e.Level, log.InfoLevel)
		assert.Equal(t, log.Fields{"file": "sloth.png", "complete": true}, e.Fields)
	}
}

func TestLogger_Trace_error(t *testing.T) {
	h := memory.New()

	l := &log.Logger{
		Handler: h,
		Level:   log.InfoLevel,
	}

	trace := l.WithField("file", "sloth.png").Trace("upload")
	err := fmt.Errorf("boom")
	trace.Stop(err)

	assert.Equal(t, 2, len(h.Entries))

	{
		e := h.Entries[0]
		assert.Equal(t, e.Message, "upload")
		assert.Equal(t, e.Level, log.InfoLevel)
		assert.Equal(t, log.Fields{"file": "sloth.png", "complete": false}, e.Fields)
	}

	{
		e := h.Entries[1]
		assert.Equal(t, e.Message, "upload")
		assert.Equal(t, e.Level, log.ErrorLevel)
		assert.Equal(t, log.Fields{"file": "sloth.png", "complete": true, "error": "boom"}, e.Fields)
	}
}

func BenchmarkLogger_small(b *testing.B) {
	l := &log.Logger{
		Handler: discard.New(),
		Level:   log.InfoLevel,
	}

	for i := 0; i < b.N; i++ {
		l.Info("login")
	}
}

func BenchmarkLogger_medium(b *testing.B) {
	l := &log.Logger{
		Handler: discard.New(),
		Level:   log.InfoLevel,
	}

	for i := 0; i < b.N; i++ {
		l.WithFields(log.Fields{
			"file": "sloth.png",
			"type": "image/png",
			"size": 1 << 20,
		}).Info("upload")
	}
}

func BenchmarkLogger_large(b *testing.B) {
	l := &log.Logger{
		Handler: discard.New(),
		Level:   log.InfoLevel,
	}

	err := fmt.Errorf("boom")

	for i := 0; i < b.N; i++ {
		l.WithFields(log.Fields{
			"file": "sloth.png",
			"type": "image/png",
			"size": 1 << 20,
		}).
			WithFields(log.Fields{
			"some":     "more",
			"data":     "here",
			"whatever": "blah blah",
			"more":     "stuff",
			"context":  "such useful",
			"much":     "fun",
		}).
			WithError(err).Error("upload failed")
	}
}
