package log_test

import (
	"context"
	"testing"

	"github.com/tj/assert"

	"github.com/apex/log"
)

func TestFromContext(t *testing.T) {
	ctx := context.Background()

	logger, ok := log.FromContext(ctx)
	assert.False(t, ok)
	assert.Nil(t, logger)

	ctx = log.NewContext(ctx, log.Log)

	logger, ok = log.FromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, log.Log, logger)
}
