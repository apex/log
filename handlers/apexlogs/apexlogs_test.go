package apexlogs_test

import (
	"testing"

	"github.com/tj/assert"

	"github.com/apex/log"
	"github.com/apex/log/handlers/apexlogs"
)

func Test(t *testing.T) {
	t.SkipNow()

	h := apexlogs.Handler{
		URL:       "http://localhost:3000",
		ProjectID: "testing",
	}

	log.SetHandler(&h)
	log.WithField("user", "tj").WithField("id", "123").Info("hello")
	log.Info("world")
	log.Error("boom")
	assert.NoError(t, h.Flush(), "flushing")
}
