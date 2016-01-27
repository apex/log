package log

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntry_WithFields(t *testing.T) {
	a := NewEntry(nil)
	assert.Nil(t, a.Fields)

	b := a.WithFields(Fields{"foo": "bar"})
	assert.Equal(t, Fields{}, a.mergedFields())
	assert.Equal(t, Fields{"foo": "bar"}, b.mergedFields())
}

func TestEntry_WithField(t *testing.T) {
	a := NewEntry(nil)
	b := a.WithField("foo", "bar")
	assert.Equal(t, Fields{}, a.mergedFields())
	assert.Equal(t, Fields{"foo": "bar"}, b.mergedFields())
}

func TestEntry_WithError(t *testing.T) {
	a := NewEntry(nil)
	b := a.WithError(fmt.Errorf("boom"))
	assert.Equal(t, Fields{}, a.mergedFields())
	assert.Equal(t, Fields{"error": "boom"}, b.mergedFields())
}
