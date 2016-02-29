package log

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLevel(t *testing.T) {
	{
		level, err := ParseLevel("info")
		assert.NoError(t, err)
		assert.Equal(t, InfoLevel, level)
	}

	{
		level, err := ParseLevel("warn")
		assert.NoError(t, err)
		assert.Equal(t, WarnLevel, level)
	}

	{
		_, err := ParseLevel("whatever")
		assert.Error(t, err)
	}
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
