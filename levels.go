package log

import (
	"bytes"
	"errors"
	"strings"
)

// Level of severity.
type Level int

// Log levels.
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var levelNames = [...]string{
	DebugLevel: "debug",
	InfoLevel:  "info",
	WarnLevel:  "warn",
	ErrorLevel: "error",
	FatalLevel: "fatal",
}

// String implementation.
func (l Level) String() string {
	return levelNames[l]
}

// MarshalJSON implementation.
func (l Level) MarshalJSON() ([]byte, error) {
	return []byte(`"` + l.String() + `"`), nil
}

// UnmarshalJSON implementation.
func (l *Level) UnmarshalJSON(b []byte) error {
	v, err := ParseLevel(string(bytes.Trim(b, `"`)))
	if err != nil {
		return err
	}

	*l = v
	return nil
}

// ParseLevel parses level string.
func ParseLevel(s string) (Level, error) {
	switch strings.ToLower(s) {
	case "debug":
		return DebugLevel, nil
	case "info":
		return InfoLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "error":
		return ErrorLevel, nil
	case "fatal":
		return FatalLevel, nil
	default:
		return -1, errors.New("invalid level")
	}
}

// MustParseLevel parses level string or panics.
func MustParseLevel(s string) Level {
	l, err := ParseLevel(s)
	if err != nil {
		panic("invalid log level")
	}

	return l
}
