package log

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLevel(t *testing.T) {
	cases := []struct {
		String string
		Level  Level
		Num    int
	}{
		{"debug", DebugLevel, 0},
		{"info", InfoLevel, 1},
		{"warn", WarnLevel, 2},
		{"warning", WarnLevel, 3},
		{"error", ErrorLevel, 4},
		{"fatal", FatalLevel, 5},
	}

	for _, c := range cases {
		t.Run(c.String, func(t *testing.T) {
			l, err := ParseLevel(c.String)
			assert.NoError(t, err, "parse")
			assert.Equal(t, c.Level, l)
		})
	}

	t.Run("invalid", func(t *testing.T) {
		l, err := ParseLevel("something")
		assert.Equal(t, ErrInvalidLevel, err)
		assert.Equal(t, InvalidLevel, l)
	})
}

func TestLevel_MarshalJSON(t *testing.T) {
	e := Entry{
		Level:   InfoLevel,
		Message: "hello",
		Fields:  Fields{},
	}

	expect := `{"fields":{},"level":"info","timestamp":"0001-01-01T00:00:00Z","message":"hello"}`

	b, err := json.Marshal(e)
	assert.NoError(t, err)
	assert.Equal(t, expect, string(b))
}

func TestLevel_UnmarshalJSON(t *testing.T) {
	s := `{"fields":{},"level":"info","timestamp":"0001-01-01T00:00:00Z","message":"hello"}`
	e := new(Entry)

	err := json.Unmarshal([]byte(s), e)
	assert.NoError(t, err)
	assert.Equal(t, InfoLevel, e.Level)
}

func TestLevel_String(t *testing.T) {
	cases := []struct {
		Level  Level
		String string
	}{
		{DebugLevel, "debug"},
		{InfoLevel, "info"},
		{WarnLevel, "warn"},
		{ErrorLevel, "error"},
		{FatalLevel, "fatal"},
		{InvalidLevel, "invalid"},
		{Level(99), "invalid"},
	}

	for _, c := range cases {
		t.Run(c.String, func(t *testing.T) {
			assert.Equal(t, c.String, c.Level.String())
		})
	}
}

func TestLevel_MarshalJSON_Invalid(t *testing.T) {
	b, err := json.Marshal(InvalidLevel)
	assert.NoError(t, err)
	assert.Equal(t, `"invalid"`, string(b))
}
