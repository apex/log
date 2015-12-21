package log

import (
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
