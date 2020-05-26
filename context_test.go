package log_test

import (
	"context"
	"testing"

	"github.com/tj/assert"

	"github.com/apex/log"
)

func TestFromContext(t *testing.T) {
	ctx := context.Background()

	logger := log.FromContext(ctx)
	assert.Equal(t, log.Log, logger)

	logs := log.WithField("foo", "bar")
	ctx = log.NewContext(ctx, logs)

	logger = log.FromContext(ctx)
	assert.Equal(t, logs, logger)
}
